import { TestBed } from '@angular/core/testing';

import { Kubernetes } from './kubernetes';

describe('Kubernetes', () => {
  let service: Kubernetes;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(Kubernetes);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
