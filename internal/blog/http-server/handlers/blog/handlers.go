package blog

import (
	"encoding/json"
	"habr/internal/blog/core/blog"
	"habr/internal/blog/http-server/dto"
	"habr/internal/pkg/formatter"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetAllBlogs(blogService *blog.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		blogs, err := blogService.GetBlogs(ctx)
		if err != nil {
			log.Println(err)
			_ = formatter.RespJSON(500, map[string]string{"error": "Internal Server Error"}, w)
		}
		err = formatter.RespJSON(200, blogs, w)
		if err != nil {
			log.Println(err)
		}
	}
}

func GetBlogByID(blogService *blog.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			_ = formatter.RespBadRequest("Invalid blog ID", w)
			return
		}
		b, err := blogService.GetBlog(ctx, id)
		if err != nil {
			log.Println(err)
			_ = formatter.RespInternalError(w)
		}
		err = formatter.RespJSON(200, b, w)
		if err != nil {
			log.Println(err)
		}
	}
}

func CreateBlog(blogService *blog.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := dto.RequestCreateBlog{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			_ = formatter.RespBadRequest("Invalid request body", w)
			return
		}

		id, err := blogService.CreateBlog(ctx, req.Name)
		if err != nil {
			log.Println(err)
			_ = formatter.RespInternalError(w)
			return
		}

		err = formatter.RespJSON(200, dto.ResponseCreateBlog{Id: id}, w)
		if err != nil {
			log.Println(err)
		}
	}
}

func UpdateBlog(blogService *blog.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid blog ID", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		req := dto.RequestUpdateBlog{}
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			_ = formatter.RespBadRequest("Invalid request body", w)
			return
		}

		id, err = blogService.UpdateBlog(ctx, req.Name, id)
		if err != nil {
			log.Println(err)
			_ = formatter.RespInternalError(w)
			return
		}

		err = formatter.RespJSON(200, dto.ResponseUpdateBlog{Id: id}, w)
		if err != nil {
			log.Println(err)
		}
	}
}

func DeleteBlog(blogService *blog.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid blog ID", http.StatusBadRequest)
			return
		}

		id, err = blogService.DeleteBlog(ctx, id)
		if err != nil {
			log.Println(err)
			_ = formatter.RespJSON(500, map[string]string{"error": "Internal Server Error"}, w)
			return
		}

		err = formatter.RespJSON(200, dto.ResponseDeleteBlog{Id: id}, w)
		if err != nil {
			log.Println(err)
		}
	}
}
