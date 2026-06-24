export interface Book {
  key: string;
  title: string;
  author_names: string[];
  first_publish_year?: number;
  cover_url?: string;
  isbn?: string[];
}

export interface SearchResponse {
  books: Book[];
  total: number;
}
