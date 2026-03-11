import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getProperties, getSearch } from '../api/client';
import type { Property, Search } from '../types';

const PAGE_SIZE = 25;

function formatPrice(price: number | null): string {
  if (price == null) return '-';
  return '$' + price.toLocaleString();
}

function formatSize(sqft: number | null): string {
  if (sqft == null) return '-';
  return sqft.toLocaleString() + ' SF';
}

export function Properties() {
  const { id } = useParams<{ id: string }>();
  const [search, setSearch] = useState<Search | null>(null);
  const [properties, setProperties] = useState<Property[]>([]);
  const [loading, setLoading] = useState(true);
  const [offset, setOffset] = useState(0);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!id) return;
    Promise.all([
      getSearch(id),
      getProperties(id, { limit: PAGE_SIZE, offset }),
    ])
      .then(([s, props]) => {
        setSearch(s);
        setProperties(props);
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'Failed to load'))
      .finally(() => setLoading(false));
  }, [id, offset]);

  if (loading) return <p>Loading...</p>;
  if (error) return <div className="error-message">{error}</div>;

  return (
    <div>
      <div className="page-header">
        <div>
          <Link to="/" className="back-link">&larr; Back to searches</Link>
          <h2>{search?.name || 'Properties'}</h2>
          {search && (
            <p className="subtitle">
              <a href={search.url} target="_blank" rel="noopener noreferrer">View on LoopNet</a>
            </p>
          )}
        </div>
      </div>

      {properties.length === 0 ? (
        <p className="empty-state">No properties found yet. Run a scrape to populate results.</p>
      ) : (
        <>
          <table className="data-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Address</th>
                <th>Type</th>
                <th>Price</th>
                <th>Size</th>
                <th>Scraped</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {properties.map((p) => (
                <tr key={p.id}>
                  <td>{p.name || '-'}</td>
                  <td>
                    {p.address}
                    {p.city && `, ${p.city}`}
                    {p.state && ` ${p.state}`}
                    {p.zip && ` ${p.zip}`}
                  </td>
                  <td>{p.propertyType || '-'}</td>
                  <td>{formatPrice(p.price)}</td>
                  <td>{formatSize(p.sizeSqFt)}</td>
                  <td>{new Date(p.scrapedAt).toLocaleDateString()}</td>
                  <td>
                    {p.url && (
                      <a href={p.url} target="_blank" rel="noopener noreferrer" className="btn btn-sm btn-secondary">
                        View
                      </a>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          <div className="pagination">
            <button
              className="btn btn-secondary btn-sm"
              disabled={offset === 0}
              onClick={() => setOffset((o) => Math.max(0, o - PAGE_SIZE))}
            >
              Previous
            </button>
            <span className="page-info">
              Showing {offset + 1}&ndash;{offset + properties.length}
            </span>
            <button
              className="btn btn-secondary btn-sm"
              disabled={properties.length < PAGE_SIZE}
              onClick={() => setOffset((o) => o + PAGE_SIZE)}
            >
              Next
            </button>
          </div>
        </>
      )}
    </div>
  );
}
