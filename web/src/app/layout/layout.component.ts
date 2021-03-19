import {
  ChangeDetectorRef,
  Component, HostListener,
  OnDestroy, OnInit,
  ViewChild
} from '@angular/core';
import { MediaMatcher } from '@angular/cdk/layout';
import { MatSidenav } from '@angular/material';

@Component({
  selector: 'app-layout',
  templateUrl: './layout.component.html',
  styleUrls: ['./layout.component.scss']
})
export class LayoutComponent implements OnInit, OnDestroy {
  @ViewChild('sidenav') sidenav: MatSidenav;
  mobileQuery: MediaQueryList;
  opened = true;

  private readonly _mobileQueryListener: () => void;
  private readonly maxWidth = 700;

  constructor(changeDetectorRef: ChangeDetectorRef, media: MediaMatcher) {
    this.mobileQuery = media.matchMedia(`(max-width: ${this.maxWidth}px)`);
    this._mobileQueryListener = () => changeDetectorRef.detectChanges();
    this.mobileQuery.addListener(this._mobileQueryListener);
  }

  ngOnInit() {
    this.checkSidenav();
  }

  ngOnDestroy() {
    this.mobileQuery.removeListener(this._mobileQueryListener);
  }

  handleClick() {
    if (window.innerWidth <= this.maxWidth && this.opened) {
      this.closeSidenav();
    }
  }

  toggle() {
    if (this.opened) {
      this.closeSidenav();
    } else {
      this.openSidenav();
    }
  }

  @HostListener('window:resize')
  onResize() {
    this.checkSidenav();
  }

  @HostListener('panright')
  openSidenav() {
    this.opened = true;
    this.sidenav.open();
  }

  @HostListener('panleft')
  closeSidenav() {
    this.opened = false;
    this.sidenav.close();
  }

  private checkSidenav() {
    if (window.innerWidth <= this.maxWidth) {
      this.closeSidenav();
    } else {
      this.openSidenav();
    }
  }
}
