import { Component, inject, signal } from '@angular/core';
import { BookService } from './book.service';
import { Book } from './models/book.model';
import { BookList } from './components/book-list/book-list';

@Component({
  selector: 'app-root',
  imports: [BookList],
  templateUrl: './app.html',
  styleUrl: './app.css',
})
export class App {
  private readonly bookService = inject(BookService);

  readonly query = signal('');
  readonly books = signal<Book[]>([]);
  readonly total = signal(0);
  readonly loading = signal(false);
  readonly error = signal('');
  readonly searched = signal(false);

  search(): void {
    const q = this.query().trim();
    if (!q) return;

    this.loading.set(true);
    this.error.set('');
    this.searched.set(true);
    this.books.set([]);

    this.bookService.searchByAuthor(q).subscribe({
      next: (res) => {
        this.books.set(res.books);
        this.total.set(res.total);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to fetch books. Make sure the backend is running on port 8080.');
        this.loading.set(false);
      },
    });
  }
}
