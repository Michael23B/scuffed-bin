import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable()
export class NotesService {
  constructor(private http: HttpClient) { }

  public listNotes() {
      return this.http.get('http://14.200.161.83:8070/post/list'); // TODO
  }

  public getNote(id) {
    return this.http.get(`http://14.200.161.83:8070/post/${id}`);
}

public postNote(text) {
      return this.http.post('http://14.200.161.83:8070/post', {'post-body': text});
  }
}