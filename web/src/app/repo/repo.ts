import { Component, computed, inject } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';

@Component({
  selector: 'app-repo',
  templateUrl: './repo.html',
  styleUrl: './repo.scss',
})
export class Repo {
  private readonly route = inject(ActivatedRoute);

  private readonly params = toSignal(this.route.paramMap, { requireSync: true });

  readonly owner = computed(() => this.params().get('owner') ?? '');
  readonly name = computed(() => this.params().get('name') ?? '');
}
