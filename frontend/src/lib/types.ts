
export interface Movie {
  id: number;
  title: string;
  description: string;
  genre: string;
  poster_image: string;
}

export interface Showtime {
  id: number;
  movie_id: number;
  start_time: string;
  capacity: number;
  reserved: number;
}
