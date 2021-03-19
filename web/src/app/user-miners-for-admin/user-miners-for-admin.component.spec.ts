import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UserMinersForAdminComponent } from './user-miners-for-admin.component';

describe('UserMinersForAdminComponent', () => {
  let component: UserMinersForAdminComponent;
  let fixture: ComponentFixture<UserMinersForAdminComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ UserMinersForAdminComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UserMinersForAdminComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
