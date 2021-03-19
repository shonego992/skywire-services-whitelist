import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UserMinerOverviewComponent } from './user-miner-overview.component';

describe('UserMinerOverviewComponent', () => {
  let component: UserMinerOverviewComponent;
  let fixture: ComponentFixture<UserMinerOverviewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ UserMinerOverviewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UserMinerOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
