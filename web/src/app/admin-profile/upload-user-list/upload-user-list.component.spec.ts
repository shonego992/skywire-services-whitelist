import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UploadUserListComponent } from './upload-user-list.component';

describe('UploadUserListComponent', () => {
  let component: UploadUserListComponent;
  let fixture: ComponentFixture<UploadUserListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ UploadUserListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UploadUserListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
