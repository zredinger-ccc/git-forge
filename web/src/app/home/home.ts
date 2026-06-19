import { Component, inject } from '@angular/core';
import { Auth } from '../core/auth';

@Component({
  selector: 'app-home',
  templateUrl: './home.html',
  styleUrl: './home.scss',
})
export class Home {
  private readonly auth = inject(Auth);

  signOut() {
    this.auth.logout().subscribe();
  }
}
