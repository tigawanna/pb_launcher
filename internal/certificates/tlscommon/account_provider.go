package tlscommon

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"pb_launcher/configs"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	privateKeyFileName  = "key.pem"
	accountJSONFileName = "account.json"
)

type LetsEncryptClientAccountProvider struct {
	baseDir string
}

func NewLetsEncryptClientAccountProvider(c configs.Config) *LetsEncryptClientAccountProvider {
	return &LetsEncryptClientAccountProvider{
		baseDir: c.GetAccountsDir(),
	}
}

func (s *LetsEncryptClientAccountProvider) buildAccountJSONFilePath(email string) string {
	return filepath.Join(s.baseDir, email, accountJSONFileName)
}

func (s *LetsEncryptClientAccountProvider) buildPrivateKeyFileName(email string) string {
	return filepath.Join(s.baseDir, email, privateKeyFileName)
}

func (s *LetsEncryptClientAccountProvider) getPrivateKey(email string) (crypto.PrivateKey, error) {
	accKeyPath := s.buildPrivateKeyFileName(email)
	if _, err := os.Stat(accKeyPath); os.IsNotExist(err) {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}

		if err := os.MkdirAll(filepath.Dir(accKeyPath), 0700); err != nil {
			return nil, err
		}
		pemBytes := certcrypto.PEMEncode(privateKey)
		if err := os.WriteFile(accKeyPath, pemBytes, 0600); err != nil {
			return nil, err
		}

		return privateKey, nil
	}

	keyBytes, err := os.ReadFile(accKeyPath)
	if err != nil {
		return nil, err
	}
	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	return x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
}

func (s *LetsEncryptClientAccountProvider) loadExistingAccount(privateKey crypto.PrivateKey, email string) (*Account, error) {
	accountFilePath := s.buildAccountJSONFilePath(email)

	fileBytes, err := os.ReadFile(accountFilePath)
	if err != nil {
		return nil, err
	}
	var account Account
	if err := json.Unmarshal(fileBytes, &account); err != nil {
		return nil, fmt.Errorf("error unmarshalling account file: %w", err)
	}
	account.Key = privateKey
	return &account, nil
}

const filePerm os.FileMode = 0o600

func (s *LetsEncryptClientAccountProvider) storeAccount(account *Account) error {
	if account == nil {
		return errors.New("account is nil")
	}
	accountFilePath := s.buildAccountJSONFilePath(account.Email)
	dir := filepath.Dir(accountFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	jsonBytes, err := json.MarshalIndent(account, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(accountFilePath, jsonBytes, filePerm)
}

func (s *LetsEncryptClientAccountProvider) SetupClient(email string) (*lego.Client, error) {
	privateKey, err := s.getPrivateKey(email)
	if err != nil {
		return nil, err
	}

	account, err := s.loadExistingAccount(privateKey, email)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if os.IsNotExist(err) {
		account = &Account{
			Email: email,
			Key:   privateKey,
		}
	}

	config := lego.NewConfig(account)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.HTTPClient = config.HTTPClient
	retryClient.Logger = nil
	config.HTTPClient = retryClient.StandardClient()

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}
	if account.Registration == nil {
		resource, err := client.Registration.Register(
			registration.RegisterOptions{TermsOfServiceAgreed: true},
		)
		if err != nil {
			return nil, err
		}
		account.Registration = resource
		if err := s.storeAccount(account); err != nil {
			return nil, err
		}
	}
	return client, err
}
