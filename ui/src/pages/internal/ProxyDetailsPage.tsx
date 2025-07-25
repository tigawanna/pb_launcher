import { Navigate, useParams, useSearchParams } from "react-router-dom";
import { useMemo, useState, type FC } from "react";
import { MenuIcon, XIcon } from "lucide-react";
import { ProxyGeneralSection } from "./proxy_section/ProxyGeneralSection";
import { DomainsSection } from "./details_section/DomainsSection";
import { useQuery, useQueryClient, type QueryKey } from "@tanstack/react-query";
import { proxyEntryService } from "../../services/proxy";

export const ProxyDetailsPage = () => {
  const { proxy_id } = useParams<{ proxy_id: string }>();
  const [searchParams, setSearchParams] = useSearchParams();
  const activeSection = searchParams.get("section") || "general";
  const [menuOpen, setMenuOpen] = useState(false);

  const rewritePathQueryKey = useMemo(
    () => ["proxy_entires", proxy_id, "rewrite_path"],
    [proxy_id],
  );
  const queryClient = useQueryClient();

  const handleSectionChange = (section: string) => {
    setSearchParams({ section });
    setMenuOpen(false);
  };

  const menuItemClass = (section: string) =>
    `btn btn-ghost justify-start w-full text-left ${activeSection === section ? "bg-primary text-primary-content" : ""}`;

  if (proxy_id == null || proxy_id === "") return <Navigate to={"/"} />;
  return (
    <div className="flex h-full flex-col md:flex-row bg-base-100 text-base-content">
      <div className="md:hidden p-4 flex justify-end items-center border-b border-base-300">
        <button
          className="btn btn-ghost"
          onClick={() => setMenuOpen(!menuOpen)}
        >
          {menuOpen ? (
            <XIcon className="w-5 h-5" />
          ) : (
            <MenuIcon className="w-5 h-5" />
          )}
        </button>
      </div>

      {(menuOpen || window.innerWidth >= 768) && (
        <aside className="w-full md:w-64 bg-base-200 p-4 border-b md:border-b-0 md:border-r border-base-300 md:block">
          <ul className="menu flex md:flex-col flex-wrap gap-2">
            <li>
              <button
                className={menuItemClass("general")}
                onClick={() => handleSectionChange("general")}
              >
                Proxy Entry
              </button>
            </li>
            <li>
              <button
                className={menuItemClass("domains")}
                onClick={() => handleSectionChange("domains")}
              >
                Domains
              </button>
            </li>
          </ul>
        </aside>
      )}

      <main className="flex-1 p-4 md:p-6 overflow-auto">
        {activeSection === "general" && (
          <div className="mb-8 ">
            <h3 className="text-lg font-semibold mb-6">General</h3>
            <div className="md:px-4 py-6 bg-base-200 rounded-box">
              <ProxyGeneralSection
                proxy_id={proxy_id}
                onChange={() =>
                  setTimeout(() =>
                    queryClient.refetchQueries({
                      queryKey: rewritePathQueryKey,
                    }),
                  )
                }
              />
            </div>
          </div>
        )}

        {activeSection === "domains" && (
          <div className="mb-8">
            <h3 className="text-lg font-semibold mb-4">Domains</h3>
            <div className="md:px-4 rounded-box">
              <DomainsSectionWrap
                queryKey={rewritePathQueryKey}
                proxy_id={proxy_id}
              />
            </div>
          </div>
        )}
      </main>
    </div>
  );
};

export const DomainsSectionWrap: FC<{
  queryKey: QueryKey;
  proxy_id: string;
}> = ({ proxy_id, queryKey }) => {
  const rewritePathQuery = useQuery({
    queryKey,
    queryFn: () => proxyEntryService.findRewritePath(proxy_id),
  });

  return (
    <DomainsSection
      proxy_id={proxy_id}
      service_id=""
      url_route_suffix={rewritePathQuery.data ?? ""}
    />
  );
};
