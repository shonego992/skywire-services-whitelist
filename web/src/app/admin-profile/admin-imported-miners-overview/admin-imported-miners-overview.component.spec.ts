import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AdminMinerOverviewComponent } from './admin-imported-miners-overview.component';

describe('AdminMinerOverviewComponent', () => {
  let component: AdminMinerOverviewComponent;
  let fixture: ComponentFixture<AdminMinerOverviewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AdminMinerOverviewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AdminMinerOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
