import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { Auth } from '../core/auth';

@Component({
  selector: 'app-login',
  imports: [FormsModule],
  templateUrl: './login.html',
  styleUrl: './login.scss',
})
export class Login {
  private readonly auth = inject(Auth);
  private readonly router = inject(Router);

  username = '';
  password = '';
  readonly error = signal<string | null>(null);
  readonly submitting = signal<boolean>(false);

  submit() {
    this.error.set(null);
    this.submitting.set(true);
    this.auth.login(this.username, this.password).subscribe({
      next: () => {
        this.submitting.set(false);
        this.router.navigateByUrl('/');
      },
      error: (err) => {
        this.submitting.set(false);
        this.error.set(this.formatError(err));
      },
    });
  }

  private formatError(err: unknown): string {
    if (err && typeof err === 'object' && 'status' in err) {
      const status = (err as { status: number }).status;
      if (status === 401) return 'Invalid username or password.';
      if (status === 0) return 'Server unreachable.';
      return `Login failed (HTTP ${status}).`;
    }
    return 'Login failed.';
  }
}
