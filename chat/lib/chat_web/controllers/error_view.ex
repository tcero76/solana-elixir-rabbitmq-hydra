defmodule ChatWeb.ErrorView do

  def render("401.json", _assigns) do
    %{error: "Unauthorized", message: "Invalid or missing token"}
  end

  def render("404.json", _assigns) do
    %{error: "Not Found"}
  end

  def render("500.json", _assigns) do
    %{error: "Internal Server Error"}
  end

  def template_not_found(_template, assigns) do
    render("500.json", assigns)
  end
end
