import { useState } from "react";
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
    return <div className="text-red-500 text-center mt-4">Error loading movies.</div>;
  }

  return (
    <div className="container mx-auto p-4">
      <header className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Movies</h1>
        {isAdmin && (
          <button
            className="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600"
            onClick={() => navigate("/movies/new")}
          >
            Add New Movie
          </button>
        )}
      </header>
      {movies.length === 0 ? (
        <p className="text-gray-600 text-center">No movies available.</p>
      ) : (
        <div className="grid grid-cols-[repeat(auto-fill,minmax(250px,1fr))] gap-8 p-8">
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
