import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup } from '@angular/forms';
import { EventResponse, KubernetesService } from '../kubernetes';
import { environment } from '../../environments/environment.development';

@Component({
  selector: 'app-gateway',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './gateway.html',
  styleUrls: ['./gateway.css']
})
export class Gateway implements OnInit {
  baseUrl = environment.baseUrl
  eventResponse?: EventResponse;
  form!: FormGroup;

  constructor(
    private kubernetesService: KubernetesService,
    private fb: FormBuilder
  ) {}

  ngOnInit(): void {
    this.form = this.fb.group({
      resourceId: ['']
    });
  }

  onClickStart(): void {
    const value = this.form.value.resourceId?.trim();
    if (!value) {
      console.warn('Resource ID is empty.');
      return;
    }

    this.kubernetesService.fetchState(value).subscribe({
      next: (response: EventResponse) => {
        this.eventResponse = response;
        console.log('Fetched event state:', response);
      },
      error: (err) => {
        console.error('Error fetching Kubernetes state:', err);
      }
    });
  }
}
