import { useState, useEffect } from 'react';

interface Movie {
  id: number;
  title: string;
  description: string;
  genre: string;
  poster_image: string;
}

interface Showtime {
  movie_id: number;
  start_time: string;
  capacity: number;
  reserved: number;
}

export function useMoviesAndShowtimes(token: string | null) {
  const [movies, setMovies] = useState<Movie[]>([]);
  const [showtimes, setShowtimes] = useState<Showtime[]>([]);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    const fetchMoviesAndShowtimes = async () => {
      try {
        if (!token) return;

        const movieResponse = await fetch('http://localhost:8080/movies', {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        });

        if (!movieResponse.ok) {
          throw new Error('Failed to fetch movies');
        }

        const moviesData = await movieResponse.json();
        setMovies(moviesData);

        const showtimeResponse = await fetch('http://localhost:8080/showtimes', {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        });

        if (!showtimeResponse.ok) {
          throw new Error('Failed to fetch showtimes');
        }

        const showtimeData = await showtimeResponse.json();
        setShowtimes(showtimeData);

      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error fetching data');
      }
    };

    fetchMoviesAndShowtimes();
  }, [token]);

  return { movies, setMovies, showtimes, setShowtimes, error };
}
