export class HttpError<T = unknown> extends Error {
  status: number;
  data?: T;

  constructor(status: number, message: string, data?: T) {
    super(message);
    this.name = "HttpError";
    this.status = status;
    this.data = data;
  }
}
