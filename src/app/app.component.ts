import { Component  } from '@angular/core';
import { NgIf  } from '@angular/common';
import { NgModel  } from '@angular/forms';
import { NotesService } from './notes.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  constructor(private notesService: NotesService) {};

  title:string = 'Scuffed bin';
  currentNote: string = '';
  currentLink: string = '';

  public listPastes(): void {
    // TODO api call to get pastes
  }

  public savePaste(): void {
    // TODO api call to save paste
    this.currentLink = 'test';
    this.notesService.listNotes().subscribe(res => console.log(res));
  }

  public clearPaste(): void {
    this.currentNote = '';
  }
}
