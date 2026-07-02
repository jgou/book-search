import { Component, input } from '@angular/core';
import { Book } from '../../models/book.model';
import { BookCard } from '../book-card/book-card';

@Component({
  selector: 'app-book-list',
  imports: [BookCard],
  templateUrl: './book-list.html',
  styleUrl: './book-list.css',
})
export class BookList {
  readonly books = input.required<Book[]>();
}
