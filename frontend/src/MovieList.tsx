import { useState, useEffect } from 'react';
import styles from './MovieList.module.css';

interface Movie {
  id: number;
  title: string;
  description: string;
  genre: string;
  poster_image: string;
}

function MovieList() {
  const [movies, setMovies] = useState<Movie[]>([]);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    const fetchMovies = async () => {
      try {
        const response = await fetch('http://localhost:8080/movies');
        if (!response.ok) {
          throw new Error('Failed to fetch movies');
        }
        const data = await response.json();
        setMovies(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error fetching movies');
      }
    };

    fetchMovies();
  }, []);

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div>
      <h1>Movies</h1>
      <div className={styles.movieGrid}>
        {movies.map((movie) => (
          <div key={movie.id} className={styles.movieCard}>
            <img src={movie.poster_image} alt={movie.title} />
            <h2>{movie.title}</h2>
            <p>{movie.description}</p>
            <p>Genre: {movie.genre}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

export default MovieList;
