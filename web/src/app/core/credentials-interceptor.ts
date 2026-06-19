import { HttpInterceptorFn } from '@angular/common/http';

// All same-origin requests carry the session cookie. PAT auth is added by
// the future tokens flow (issue #5) when the UI grows a token-management
// surface.
export const credentialsInterceptor: HttpInterceptorFn = (req, next) => {
  return next(req.clone({ withCredentials: true }));
};
