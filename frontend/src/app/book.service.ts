import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { SearchResponse } from './models/book.model';

@Injectable({ providedIn: 'root' })
export class BookService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = 'http://localhost:8080/api';

  searchByAuthor(author: string, limit = 20): Observable<SearchResponse> {
    const params = new HttpParams()
      .set('author', author)
      .set('limit', String(limit));
    return this.http.get<SearchResponse>(`${this.apiUrl}/books`, { params });
  }
}
