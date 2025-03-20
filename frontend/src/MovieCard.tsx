"use client";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Plus, Trash } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./components/ui/dialog";
import SeatSelection from "./SeatSelection";
import { useState } from "react";

interface Showtime {
  id: number;
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

function MovieCard({
  movie,
  showtimes,
  isAdmin,
  onDelete,
  onAddShowtime,
}: MovieCardProps) {
  const [selectedSeats, setSelectedSeats] = useState<string[]>([]);

  const movieShowtimes =
    showtimes?.filter((showtime) => showtime.movie_id === movie.id) || [];
  const groupedShowtimes = groupShowtimesByDate(movieShowtimes);

  const handleGetSeats = async (showtimeId: number) => {
    const token = localStorage.getItem("token");
    try {
      const response = await fetch(
        `http://localhost:8080/showtimes/seats/${showtimeId}`,
        {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        }
      );

      if (!response.ok) {
        throw new Error(`Error fetching seats: ${response.statusText}`);
      }

      const availableSeats = await response.json();
      setSelectedSeats(availableSeats);
      console.log(selectedSeats);
    } catch (error) {
      console.error("Failed to fetch seats:", error);
      setSelectedSeats([]);
    }
  };

  return (
    <Card className="overflow-hidden p-0">
      <div
        className="relative w-full h-[375px] bg-cover bg-center"
        style={{
          backgroundImage: `url(${
            movie.poster_image || "/placeholder.svg?height=375&width=250"
          })`,
        }}
      ></div>

      <CardHeader className="flex flex-row items-center justify-between px-4 space-y-0">
        <h2 className="text-xl font-semibold">{movie.title}</h2>
        {isAdmin && (
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onDelete(movie.id)}
          >
            <Trash className="h-4 w-4" />
          </Button>
        )}
      </CardHeader>

      <CardContent className="px-4 pb-2 space-y-2">
        <p className="text-sm text-muted-foreground">{movie.description}</p>
        <p className="text-sm text-muted-foreground italic">
          Genre: {movie.genre}
        </p>
      </CardContent>

      <Separator />

      <CardFooter className="flex flex-col items-start p-4">
        <div className="w-full">
          {isAdmin && (
            <Button
              variant="ghost"
              size="sm"
              className="px-4 py-2 h-auto mb-2 font-normal"
              onClick={() => onAddShowtime(movie.id)}
            >
              <Plus className="h-3.5 w-3.5 mr-1" />
              Add Showtime
            </Button>
          )}

          <h4 className="font-medium mb-2">Showtimes:</h4>

          {Object.keys(groupedShowtimes).length > 0 ? (
            Object.entries(groupedShowtimes).map(([date, dateShowtimes]) => (
              <div key={date} className="mb-4 w-full">
                <h5 className="font-medium text-sm mb-2">{date}</h5>
                <div className="flex flex-wrap gap-2">
                  {dateShowtimes.map((showtime, index) => (
                    <Dialog>
                      <DialogTrigger onClick={() => handleGetSeats(showtime.id)} asChild>
                        <Button
                          key={index}
                          variant="secondary"
                          size="sm"
                          className="text-xs"
                        >
                          {formatTime(showtime.start_time)}
                        </Button>
                      </DialogTrigger>
                      <DialogContent>
                        <DialogHeader>
                          <DialogTitle>
                            {movie.title} at {formatTime(showtime.start_time)}{" "}
                            {date}
                          </DialogTitle>
                        </DialogHeader>
                        <SeatSelection availableSeats={selectedSeats} />
                      </DialogContent>
                    </Dialog>
                  ))}
                </div>
              </div>
            ))
          ) : (
            <p className="text-muted-foreground italic">
              No showtimes available
            </p>
          )}
        </div>
      </CardFooter>
    </Card>
  );
}

export default MovieCard;
