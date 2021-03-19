import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { WhitelistAppComponent } from './whitelist-app.component';

describe('WhitelistAppComponent', () => {
  let component: WhitelistAppComponent;
  let fixture: ComponentFixture<WhitelistAppComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ WhitelistAppComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WhitelistAppComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
