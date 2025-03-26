"use client";

import type React from "react";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Movie } from "./lib/types";

interface NewMovieFormProps {
  onSuccess: (movie: Movie) => void;
  onCancel: () => void;
}

export function NewMovieForm({ onSuccess, onCancel }: NewMovieFormProps) {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [genre, setGenre] = useState("");
  const [posterImage, setPosterImage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError("");

    try {
      const token = localStorage.getItem("token");
      if (!token) {
        throw new Error("No authentication token found");
      }

      const response = await fetch("http://localhost:8080/movies/add", {
        method: "POST",
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
      });

      if (!response.ok) {
        throw new Error("Failed to add movie");
      }

      const newMovie = await response.json();
      onSuccess({
        id: newMovie.id,
        title: newMovie.title,
        description: newMovie.description,
        genre: newMovie.genre,
        poster_image: newMovie.poster_image,
      });
    } catch (error) {
      console.error("Error adding movie:", error);
      setError(error instanceof Error ? error.message : "An error occurred");
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="bg-destructive/15 text-destructive text-sm p-3 rounded-md">
          {error}
        </div>
      )}

      <div className="space-y-2">
        <Label htmlFor="title">Title</Label>
        <Input
          id="title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="description">Description</Label>
        <Textarea
          id="description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          className="min-h-[120px]"
          required
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="genre">Genre</Label>
        <Input
          id="genre"
          type="text"
          value={genre}
          onChange={(e) => setGenre(e.target.value)}
          required
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="posterImage">Poster Image URL</Label>
        <Input
          id="posterImage"
          type="url"
          value={posterImage}
          onChange={(e) => setPosterImage(e.target.value)}
          placeholder="https://example.com/movie-poster.jpg"
          required
        />
      </div>

      <div className="flex justify-end gap-2 pt-2">
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={isSubmitting}
        >
          Cancel
        </Button>
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? "Adding Movie..." : "Add Movie"}
        </Button>
      </div>
    </form>
  );
}
