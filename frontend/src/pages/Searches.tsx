import { useState, useEffect, type FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { getSearches, createSearch, updateSearch, deleteSearch, triggerScrape } from '../api/client';
import type { Search, SearchFormData } from '../types';

const emptyForm: SearchFormData = {
  name: '',
  url: '',
  active: true,
};

export function Searches() {
  const navigate = useNavigate();
  const [searches, setSearches] = useState<Search[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState<SearchFormData>(emptyForm);
  const [error, setError] = useState('');

  useEffect(() => {
    loadSearches();
  }, []);

  async function loadSearches() {
    try {
      const data = await getSearches();
      setSearches(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load searches');
    } finally {
      setLoading(false);
    }
  }

  function openCreate() {
    setForm(emptyForm);
    setEditingId(null);
    setShowForm(true);
  }

  function openEdit(s: Search) {
    setForm({
      name: s.name,
      url: s.url,
      active: s.active,
    });
    setEditingId(s.id);
    setShowForm(true);
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    try {
      if (editingId) {
        await updateSearch(editingId, form);
      } else {
        await createSearch(form);
      }
      setShowForm(false);
      setEditingId(null);
      await loadSearches();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save search');
    }
  }

  async function handleDelete(id: string) {
    if (!confirm('Delete this search and all its properties?')) return;
    try {
      await deleteSearch(id);
      await loadSearches();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete search');
    }
  }

  async function handleScrape(id: string) {
    try {
      await triggerScrape(id);
      alert('Scrape started');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to trigger scrape');
    }
  }

  function setField<K extends keyof SearchFormData>(key: K, value: SearchFormData[K]) {
    setForm((prev) => ({ ...prev, [key]: value }));
  }

  if (loading) return <p>Loading...</p>;

  return (
    <div>
      <div className="page-header">
        <h2>Saved Searches</h2>
        <button onClick={openCreate} className="btn btn-primary">
          New Search
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}

      {showForm && (
        <div className="modal-backdrop" onClick={() => setShowForm(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>{editingId ? 'Edit Search' : 'New Search'}</h3>
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Name</label>
                <input value={form.name} onChange={(e) => setField('name', e.target.value)} required />
              </div>
              <div className="form-group">
                <label>LoopNet URL</label>
                <input
                  value={form.url}
                  onChange={(e) => setField('url', e.target.value)}
                  placeholder="https://www.loopnet.com/search/..."
                  required
                />
              </div>
              <div className="form-group">
                <label className="checkbox-label">
                  <input
                    type="checkbox"
                    checked={form.active}
                    onChange={(e) => setField('active', e.target.checked)}
                  />
                  Active
                </label>
              </div>
              <div className="form-actions">
                <button type="submit" className="btn btn-primary">
                  {editingId ? 'Update' : 'Create'}
                </button>
                <button type="button" className="btn btn-secondary" onClick={() => setShowForm(false)}>
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {searches.length === 0 ? (
        <p className="empty-state">No searches configured yet. Create one to get started.</p>
      ) : (
        <table className="data-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>URL</th>
              <th>Status</th>
              <th>Created</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {searches.map((s) => (
              <tr key={s.id}>
                <td>
                  <a href="#" onClick={(e) => { e.preventDefault(); navigate(`/searches/${s.id}`); }}>
                    {s.name}
                  </a>
                </td>
                <td>
                  <a href={s.url} target="_blank" rel="noopener noreferrer">
                    {s.url.length > 60 ? s.url.slice(0, 60) + '…' : s.url}
                  </a>
                </td>
                <td>
                  <span className={`badge ${s.active ? 'badge-active' : 'badge-inactive'}`}>
                    {s.active ? 'Active' : 'Inactive'}
                  </span>
                </td>
                <td>{new Date(s.createdAt).toLocaleDateString()}</td>
                <td className="actions">
                  <button onClick={() => openEdit(s)} className="btn btn-sm btn-secondary">Edit</button>
                  <button onClick={() => handleScrape(s.id)} className="btn btn-sm btn-primary">Scrape</button>
                  <button onClick={() => handleDelete(s.id)} className="btn btn-sm btn-danger">Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
