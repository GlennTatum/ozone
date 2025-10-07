import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map } from 'rxjs';
import { environment } from '../environments/environment.development';

export interface EventResponse {
  PortAssignment: number;
}

@Injectable({
  providedIn: 'root'
})
export class KubernetesService {
  private state?: EventResponse;

  constructor(private http: HttpClient) {}

  /**
   * Calls the API endpoint, extracts PortAssignment,
   * stores it locally, and returns it as an Observable<number>.
   */
  fetchState(id: string): Observable<EventResponse> {
    return this.http.post<EventResponse>(`${environment.api}/events/e/${id}`, { withCredentials: true }).pipe(
      map((response: EventResponse) => {
        this.state = response;
        return this.state;
      })
    );
  }

  /**
   * Returns the last cached port value (if any).
   */
  getPort(): number | undefined {
    return this.state?.PortAssignment;
  }
}
