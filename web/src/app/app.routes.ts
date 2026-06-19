import { Routes } from '@angular/router';
import { authGuard } from './core/auth-guard';

export const routes: Routes = [
  {
    path: 'login',
    loadComponent: () => import('./login/login').then((m) => m.Login),
  },
  {
    path: '',
    canActivate: [authGuard],
    loadComponent: () => import('./home/home').then((m) => m.Home),
    pathMatch: 'full',
  },
  {
    path: 'r/:owner/:name',
    canActivate: [authGuard],
    loadComponent: () => import('./repo/repo').then((m) => m.Repo),
  },
  { path: '**', redirectTo: '' },
];
