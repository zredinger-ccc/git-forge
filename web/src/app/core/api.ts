import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Health {
  status: string;
  build: { version: string; commit: string; date: string };
}

@Injectable({ providedIn: 'root' })
export class Api {
  private readonly http = inject(HttpClient);

  health(): Observable<Health> {
    return this.http.get<Health>('/healthz');
  }

  login(username: string, password: string): Observable<void> {
    return this.http.post<void>('/api/auth/login', { username, password });
  }

  logout(): Observable<void> {
    return this.http.post<void>('/api/auth/logout', {});
  }
}
