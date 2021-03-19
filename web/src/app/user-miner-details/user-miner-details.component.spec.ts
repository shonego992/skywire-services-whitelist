import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UserMinerDetailsComponent } from './user-miner-details.component';

describe('UserMinerDetailsComponent', () => {
  let component: UserMinerDetailsComponent;
  let fixture: ComponentFixture<UserMinerDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ UserMinerDetailsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UserMinerDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
