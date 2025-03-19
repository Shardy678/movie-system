import { useState } from 'react';
import styles from './MovieList.module.css';
import { useAuth } from './useAuth';
import { useMoviesAndShowtimes } from './useMoviesAndShowtimes';
import { useNavigate } from 'react-router-dom';
import ShowtimeForm from './ShowtimeForm';

interface Showtime {
  movie_id: number;
  start_time: string;
  capacity: number;
  reserved: number;
}

interface GroupedShowtimes {
  [date: string]: Showtime[];
}

function MovieList() {
  const [showForm, setShowForm] = useState<boolean>(false);
  const [newShowtime, setNewShowtime] = useState({ movie_id: 0, start_time: '', capacity: 0 });
  const [selectedMovieId, setSelectedMovieId] = useState<number | null>(null);
  const { isAdmin, token } = useAuth();
  const { movies, showtimes, setShowtimes, setMovies, error } = useMoviesAndShowtimes(token);
  const navigate = useNavigate(); 


  async function handleDelete(movieID: number) {
    if (!token) {
      console.error('No token found, user might not be authenticated.');
      return;
    }

    try {
      const response = await fetch(`http://localhost:8080/movies/delete/${movieID}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        throw new Error('Failed to delete the movie');
      }

      setMovies(prevMovies => prevMovies.filter(movie => movie.id !== movieID));
    } catch (error) {
      console.error('Error deleting movie:', error);
    }
  }

  async function handleAddShowtime(e: React.FormEvent) {
    e.preventDefault();
    
    if (!newShowtime.start_time || newShowtime.capacity <= 0) {
      alert('Please enter a valid showtime and capacity.');
      return;
    }

    if (!token) {
      console.error('No token found, user might not be authenticated.');
      return;
    }

    try {
      const response = await fetch('http://localhost:8080/showtimes/add', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(newShowtime)
      });

      if (!response.ok) {
        throw new Error('Failed to add showtime');
      }

      const addedShowtime = await response.json();
      setShowtimes([...showtimes, addedShowtime]);
      setShowForm(false);
      setNewShowtime({ movie_id: 0, start_time: '', capacity: 0 });
    } catch (error) {
      console.error('Error adding showtime:', error);
    }
  }

  const formatDate = (dateString: string): string =>
    new Date(dateString).toLocaleDateString([], { weekday: 'long', month: 'long', day: 'numeric' });

  const formatTime = (dateString: string): string =>
    new Date(dateString).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });

  const groupShowtimesByDate = (showtimes: Showtime[]): GroupedShowtimes =>
    showtimes.reduce((acc, showtime) => {
      const date = formatDate(showtime.start_time);
      acc[date] = acc[date] || [];
      acc[date].push(showtime);
      return acc;
    }, {} as GroupedShowtimes);

  if (error) {
    return <div>Error</div>
  }

  return (
    <div className={styles.container}>
       <header>
        <h1>Movies</h1>

        {isAdmin && (
          <button className={styles.addButton} onClick={() => navigate('/movies/new')}>
            Add New Movie
          </button>
        )}
      </header>

      {movies.length === 0 ? (
        <p className={styles.noMovies}>No movies available.</p>
      ) : (
        <div className={styles.movieGrid}>
          {movies.map((movie) => {
            const movieShowtimes = showtimes.filter(showtime => showtime.movie_id === movie.id);
            const groupedShowtimes = groupShowtimesByDate(movieShowtimes);

            return (
              <div key={movie.id} className={styles.movieCard}>
                <img src={movie.poster_image || "/placeholder.svg"} alt={movie.title} />
                <div className={styles.movieTitle}>
                  <h2>{movie.title}</h2>
                  {isAdmin && (
                    <button className={styles.deleteButton} onClick={() => handleDelete(movie.id)}>
                      <i className="fas fa-trash"></i>
                    </button>
                  )}
                </div>
                <p>{movie.description}</p>
                <p>Genre: {movie.genre}</p>

                <div className={styles.showtimes}>
                  <h4 className={styles.showtimeHeading}>Showtimes:</h4>
                  {Object.keys(groupedShowtimes).length > 0 ? (
                    Object.entries(groupedShowtimes).map(([date, dateShowtimes]) => (
                      <div key={date} className={styles.showtimeDate}>
                        <h5>{date}</h5>
                        <div className={styles.showtimeContainer}>
                          {dateShowtimes.map((showtime, index) => (
                            <button key={index} className={styles.showtimeButton}>
                              {formatTime(showtime.start_time)}
                            </button>
                          ))}
                        </div>
                      </div>
                    ))
                  ) : (
                    <p className={styles.noShowtimes}>No showtimes available</p>
                  )}
                </div>

                {isAdmin && (
                  <button className={styles.addShowtimeButton} onClick={() => {
                  setSelectedMovieId(movie.id);
                  setShowForm(true);
                }}>
                  + Add Showtime
                </button>
                )}
              </div>
            );
          })}
        </div>
      )}

        {showForm && selectedMovieId !== null && (
        <ShowtimeForm
          movieId={selectedMovieId}
          onClose={() => setShowForm(false)}
          onAddShowtime={() => handleAddShowtime}
        />
      )}    </div>
  );
}

export default MovieList;
