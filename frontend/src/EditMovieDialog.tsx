import { Pencil } from "lucide-react";
import { Button } from "./components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./components/ui/dialog";
import { Movie } from "./lib/types";
import React, { useState } from "react";
import { Label } from "./components/ui/label";
import { Input } from "./components/ui/input";
import { Textarea } from "./components/ui/textarea";

interface EditMovieDialogProps {
  movie: Movie;
  onMovieUpdate: (updatedMovie: Movie) => void;
}

const EditMovieDialog = ({ movie, onMovieUpdate }: EditMovieDialogProps) => {
  const [open, setOpen] = useState(false);
  const [title, setTitle] = useState(movie.title);
  const [description, setDescription] = useState(movie.description);
  const [genre, setGenre] = useState(movie.genre);
  const [posterImage, setPosterImage] = useState(movie.poster_image || "");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    const token = localStorage.getItem("token");

    try {
      const response = await fetch(
        `http://localhost:8080/movies/update/${movie.id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            title,
            description,
            genre,
            poster_image: posterImage,
          }),
        }
      );

      if (!response.ok) {
        throw new Error(`Error updating movie: ${response.statusText}`);
      }
      const updatedMovie = await response.json();
      onMovieUpdate(updatedMovie);
      setOpen(false);
    } catch (err: any) {
      setError(err.message || "Something went wrong");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="ghost" size="icon">
          <Pencil />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Movie</DialogTitle>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <Label htmlFor="title" className="mb-1">
                Title
              </Label>
              <Input
                className="mt-1"
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
              />
            </div>
            <div>
              <Label className="mb-1" htmlFor="description">
                Description
              </Label>
              <Textarea
                className="mt-1"
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </div>
            <div>
              <Label className="mb-1" htmlFor="genre">
                Genre
              </Label>
              <Input
                className="mt-1"
                id="genre"
                value={genre}
                onChange={(e) => setGenre(e.target.value)}
              />
            </div>
            <div>
              <Label className="mb-1" htmlFor="posterImage">
                Poster Image URL
              </Label>
              <Input
                className="mt-1"
                id="posterImage"
                value={posterImage}
                onChange={(e) => setPosterImage(e.target.value)}
              />
            </div>
            {error && <p className="text-red-500">{error}</p>}
            <Button type="submit" disabled={loading}>
              {loading ? "Updating..." : "Update Movie"}
            </Button>
          </form>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
};

export default EditMovieDialog;
