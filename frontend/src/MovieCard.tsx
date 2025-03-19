import styles from "./MovieList.module.css";

interface Showtime {
  movie_id: number;
  start_time: string;
  capacity: number;
  reserved: number;
}

interface GroupedShowtimes {
  [date: string]: Showtime[];
}

interface Movie {
  id: number;
  title: string;
  description: string;
  genre: string;
  poster_image?: string;
}

interface MovieCardProps {
  movie: Movie;
  showtimes: Showtime[];
  isAdmin: boolean;
  onDelete: (movieId: number) => void;
  onAddShowtime: (movieId: number) => void;
}

// Utility functions for formatting
const formatDate = (dateString: string): string =>
  new Date(dateString).toLocaleDateString([], {
    weekday: "long",
    month: "long",
    day: "numeric",
  });

const formatTime = (dateString: string): string =>
  new Date(dateString).toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });

const groupShowtimesByDate = (showtimes: Showtime[]): GroupedShowtimes =>
  showtimes.reduce((acc, showtime) => {
    const date = formatDate(showtime.start_time);
    acc[date] = acc[date] || [];
    acc[date].push(showtime);
    return acc;
  }, {} as GroupedShowtimes);

function MovieCard({ movie, showtimes, isAdmin, onDelete, onAddShowtime }: MovieCardProps) {
  const movieShowtimes = showtimes.filter((showtime) => showtime.movie_id === movie.id);
  const groupedShowtimes = groupShowtimesByDate(movieShowtimes);

  return (
    <div className={styles.movieCard}>
      <img src={movie.poster_image || "/placeholder.svg"} alt={movie.title} />
      <div className={styles.movieTitle}>
        <h2>{movie.title}</h2>
        {isAdmin && (
          <button className={styles.deleteButton} onClick={() => onDelete(movie.id)}>
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
        <button className={styles.addShowtimeButton} onClick={() => onAddShowtime(movie.id)}>
          + Add Showtime
        </button>
      )}
    </div>
  );
}

export default MovieCard;
