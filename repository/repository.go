package repository


func (r *Repository) Close() error {
    return r.db.Close()
}