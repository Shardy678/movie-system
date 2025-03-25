"use client";

import { useState } from "react";
import { useAuth } from "./useAuth";
import { useMoviesAndShowtimes } from "./useMoviesAndShowtimes";
import ShowtimeForm from "./ShowtimeForm";
import MovieCard from "./MovieCard";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { NewMovieForm } from "./NewMovieForm";
import { Movie } from "./lib/types";

function MovieList() {
  const [showForm, setShowForm] = useState<boolean>(false);
  const [showNewMovieDialog, setShowNewMovieDialog] = useState<boolean>(false);
  const [selectedMovieId, setSelectedMovieId] = useState<number | null>(null);
  const { isAdmin, token } = useAuth();
  const { movies, showtimes, setMovies, error } = useMoviesAndShowtimes(token);

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

  const handleAddMovie = (newMovie: Movie) => {
    console.log("New movie added:", newMovie);
    setMovies((prev) => [...prev, newMovie]);
    setShowNewMovieDialog(false);
  };

  if (error) {
    return (
      <div className="text-red-500 text-center mt-4">Error loading movies.</div>
    );
  }

  return (
    <div className="container mx-auto p-4">
      <header className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Movies</h1>
        {isAdmin && (
          <Button onClick={() => setShowNewMovieDialog(true)}>
            Add New Movie
          </Button>
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
        />
      )}

      <Dialog open={showNewMovieDialog} onOpenChange={setShowNewMovieDialog}>
        <DialogContent className="sm:max-w-[550px]">
          <DialogHeader>
            <DialogTitle>Add New Movie</DialogTitle>
            <DialogDescription>
              Enter the details of the movie you want to add to the database.
            </DialogDescription>
          </DialogHeader>
          <NewMovieForm
            onSuccess={handleAddMovie}
            onCancel={() => setShowNewMovieDialog(false)}
          />
        </DialogContent>
      </Dialog>
    </div>
  );
}

export default MovieList;
