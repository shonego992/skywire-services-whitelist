import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AdminAllMinersViewComponent } from './admin-all-miners-view.component';

describe('AdminAllMinersViewComponent', () => {
  let component: AdminAllMinersViewComponent;
  let fixture: ComponentFixture<AdminAllMinersViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AdminAllMinersViewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AdminAllMinersViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
