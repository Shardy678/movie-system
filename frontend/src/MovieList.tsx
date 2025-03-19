import { useState } from "react";
import styles from "./MovieList.module.css";
import { useAuth } from "./useAuth";
import { useMoviesAndShowtimes } from "./useMoviesAndShowtimes";
import { useNavigate } from "react-router-dom";
import ShowtimeForm from "./ShowtimeForm";
import MovieCard from "./MovieCard";


function MovieList() {
  const [showForm, setShowForm] = useState<boolean>(false);
  const [selectedMovieId, setSelectedMovieId] = useState<number | null>(null);
  const { isAdmin, token } = useAuth();
  const { movies, showtimes, setShowtimes, setMovies, error } =
    useMoviesAndShowtimes(token);
  const navigate = useNavigate();

  async function deleteMovie(movieId: number) {
    if (!token) {
      console.error("No token found, user might not be authenticated.");
      return;
    }
    try {
      const response = await fetch(
        `http://localhost:8080/movies/delete/${movieId}`,
        {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        }
      );
      if (!response.ok) {
        throw new Error("Failed to delete the movie");
      }
      setMovies((prev) => prev.filter((movie) => movie.id !== movieId));
    } catch (error) {
      console.error("Error deleting movie:", error);
    }
  }

  async function addShowtime(newShowtime: {
    movie_id: number;
    start_time: string;
    capacity: number;
  }) {
    console.log("New Showtime object:", newShowtime);
    if (!token) {
      console.error("No token found, user might not be authenticated.");
      return;
    }
    try {
      const response = await fetch("http://localhost:8080/showtimes/add", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(newShowtime),
      });
  
      if (!response.ok) {
        throw new Error("Failed to add showtime");
      }
  
      const addedShowtime = await response.json();
  
      setShowtimes((prevShowtimes) => [
        ...(Array.isArray(prevShowtimes) ? prevShowtimes : []),
        addedShowtime,
      ]);  
      setShowForm(false);
    } catch (error) {
      console.error("Error adding showtime:", error);
    }
  }
  

  if (error) {
    return <div>Error</div>;
  }

  return (
    <div className={styles.container}>
      <header>
        <h1>Movies</h1>
        {isAdmin && (
          <button
            className={styles.addButton}
            onClick={() => navigate("/movies/new")}
          >
            Add New Movie
          </button>
        )}
      </header>
      {movies.length === 0 ? (
        <p className={styles.noMovies}>No movies available.</p>
      ) : (
        <div className={styles.movieGrid}>
          {movies.map((movie) => (
            <MovieCard
              key={movie.id}
              movie={movie}
              showtimes={showtimes}
              isAdmin={isAdmin}
              onDelete={deleteMovie}
              onAddShowtime={(id: number) => {
                setSelectedMovieId(id);
                setShowForm(true);
              }}
            />
          ))}
        </div>
      )}
      {showForm && selectedMovieId !== null && (
        <ShowtimeForm
          movieId={selectedMovieId}
          onClose={() => setShowForm(false)}
          onAddShowtime={addShowtime}
        />
      )}
    </div>
  );
}

export default MovieList;
