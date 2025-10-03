import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Observable, timer } from 'rxjs';

@Component({
  selector: 'app-gateway',
  imports: [],
  templateUrl: './gateway.html',
  styleUrl: './gateway.css'
})
export class Gateway implements OnInit {

  id: string
  discovery: boolean

  constructor(private activatedRoute: ActivatedRoute) {
    this.id = ""
    this.discovery = false
  }

  public ngOnInit(): void {
      this.activatedRoute.url.subscribe((segment) => {
        const p = segment.at(1)
        this.id = p ? p.toString() : ""
      })
  }

  public Discover() {
    // setup rxjs timer to poll api server for pod availability
    this.discovery = true
  }
}
