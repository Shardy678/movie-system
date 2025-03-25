import { useState, useEffect } from "react";
import { Movie } from "./lib/types";

interface Showtime {
  id: number;
  movie_id: number;
  start_time: string;
  capacity: number;
  reserved: number;
}

export function useMoviesAndShowtimes(token: string | null) {
  const [movies, setMovies] = useState<Movie[]>([]);
  const [showtimes, setShowtimes] = useState<Showtime[]>([]);
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    const fetchMoviesAndShowtimes = async () => {
      if (!token) return;

      setLoading(true);
      try {
        const movieResponse = await fetch("http://localhost:8080/movies", {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        });

        if (!movieResponse.ok) {
          throw new Error("Failed to fetch movies");
        }

        const moviesData = await movieResponse.json();
        setMovies(moviesData);

        const showtimeResponse = await fetch(
          "http://localhost:8080/showtimes",
          {
            headers: {
              Authorization: `Bearer ${token}`,
              "Content-Type": "application/json",
            },
          }
        );

        if (!showtimeResponse.ok) {
          throw new Error("Failed to fetch showtimes");
        }

        const showtimeData = await showtimeResponse.json();
        const uniqueShowtimes = showtimeData.filter(
          (showtime: Showtime, index: number, self: Showtime[]) =>
            index === self.findIndex((s) => s.id === showtime.id)
        );
        setShowtimes(uniqueShowtimes);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Error fetching data");
        console.error("Error fetching data:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchMoviesAndShowtimes();

    return () => {
      setMovies([]);
      setShowtimes([]);
      setError("");
    };
  }, [token]);

  return { movies, setMovies, showtimes, setShowtimes, error, loading };
}
