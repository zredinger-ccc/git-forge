import { Injectable, signal, inject } from '@angular/core';
import { Router } from '@angular/router';
import { tap } from 'rxjs/operators';
import { Api } from './api';

const STORAGE_KEY = 'git-forge.session';

@Injectable({ providedIn: 'root' })
export class Auth {
  private readonly api = inject(Api);
  private readonly router = inject(Router);

  // Coarse signal — true if we believe a session cookie is set.
  // The real source of truth is the server; this just gates the UI.
  readonly authenticated = signal<boolean>(this.readPersisted());

  login(username: string, password: string) {
    return this.api.login(username, password).pipe(
      tap(() => {
        this.setAuthenticated(true);
      }),
    );
  }

  logout() {
    return this.api.logout().pipe(
      tap(() => {
        this.setAuthenticated(false);
        this.router.navigateByUrl('/login');
      }),
    );
  }

  private setAuthenticated(value: boolean) {
    this.authenticated.set(value);
    if (value) {
      localStorage.setItem(STORAGE_KEY, '1');
    } else {
      localStorage.removeItem(STORAGE_KEY);
    }
  }

  private readPersisted(): boolean {
    try {
      return localStorage.getItem(STORAGE_KEY) === '1';
    } catch {
      return false;
    }
  }
}
