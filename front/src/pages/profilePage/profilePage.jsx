import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import classes from "./profile.module.css";
import LoadingSpinner from "../../components/loadingSpinner/loadingSpinner";

const ProfilePage = () => {
  const navigate = useNavigate();
  const API = 'http://localhost:8080';

  const [user, setUser] = useState({ name: '', email: '' });
  const [settings, setSettings] = useState({ timezone: 'UTC' });
  const [isEditing, setIsEditing] = useState(false);
  const [editForm, setEditForm] = useState({ name: '', email: '', password: '', timezone: 'UTC' });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const timezoneOptions = [
    'UTC-12', 'UTC-11', 'UTC-10', 'UTC-9', 'UTC-8', 'UTC-7', 'UTC-6', 'UTC-5',
    'UTC-4', 'UTC-3', 'UTC-2', 'UTC-1', 'UTC', 'UTC+1', 'UTC+2', 'UTC+3',
    'UTC+4', 'UTC+5', 'UTC+6', 'UTC+7', 'UTC+8', 'UTC+9', 'UTC+10', 'UTC+11', 'UTC+12'
  ];

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);

        const userRes = await fetch(`${API}/v1/private/users`, {
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include'
        });
        if (!userRes.ok) throw new Error('Failed to load user');
        const userData = await userRes.json();

        const settingsRes = await fetch(`${API}/v1/private/settings`, {
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include'
        });
        if (!settingsRes.ok) throw new Error('Failed to load settings');
        const settingsData = await settingsRes.json();

        const u = userData.user || userData;
        const s = settingsData.settings || settingsData;

        setUser({ name: u.user_name || '', email: u.user_email || '' });
        setSettings({ timezone: s.user_tz || 'UTC' });
      } catch (e) {
        setError(e.message);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const handleEditClick = () => {
    setEditForm({
      name: user.name,
      email: user.email,
      password: '',
      timezone: settings.timezone
    });
    setIsEditing(true);
  };

  const handleCancelEdit = () => {
    setIsEditing(false);
  };

  const handleUpdateUser = async (e) => {
    e.preventDefault();

    const userChanged = editForm.name !== user.name || editForm.email !== user.email || editForm.password !== '';
    const tzChanged = editForm.timezone !== settings.timezone;

    if (!userChanged && !tzChanged) {
      setIsEditing(false);
      return;
    }

    try {
      if (userChanged) {
        const userBody = {};
        if (editForm.name !== user.name) userBody.user_name = editForm.name;
        if (editForm.email !== user.email) userBody.user_email = editForm.email;
        if (editForm.password) userBody.user_pass = editForm.password;

        const res = await fetch(`${API}/v1/private/users`, {
          method: 'PATCH',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify(userBody)
        });
        if (!res.ok) {
          const err = await res.text();
          throw new Error(err || 'User update failed');
        }

        setUser({ name: editForm.name, email: editForm.email });
      }

      if (tzChanged) {
        const res = await fetch(`${API}/v1/private/settings`, {
          method: 'PATCH',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({ user_tz: editForm.timezone })
        });
        if (!res.ok) {
          const err = await res.text();
          throw new Error(err || 'Settings update failed');
        }

        setSettings({ timezone: editForm.timezone });
      }

      setIsEditing(false);
    } catch (err) {
      alert('Failed to update: ' + err.message);
    }
  };

  const handleDeleteUser = async () => {
    if (!window.confirm('Delete your account? This cannot be undone.')) return;
    const res = await fetch(`${API}/v1/private/users`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include'
    });
    if (res.ok) navigate('/auth');
    else alert('Failed to delete user');
  };

  if (loading) return (
    <div className={classes.loadingContainer}>
      <LoadingSpinner />
      <div className={classes.loadingText}>Loading profile...</div>
    </div>
  );

  if (error) return (
    <div className={classes.errorContainer}>
      <div className={classes.errorText}>Error: {error}</div>
      <button className={classes.retryButton} onClick={() => window.location.reload()}>Try again</button>
    </div>
  );

  return (
    <div className={classes.userProfile}>
      <div className={classes.profileContainer}>

        <div className={classes.profileHeader}>
          <div className={classes.avatarWrapper}>
            <div className={classes.avatarRing}></div>
            <div className={classes.avatar}>
              {user.name ? user.name.charAt(0).toUpperCase() : 'U'}
            </div>
          </div>
          <h2 className={classes.title}>{user.name}</h2>
          <span className={classes.subtitle}>{user.email}</span>
        </div>

        {!isEditing ? (
          <>
            <div className={classes.timeline}>
              <div className={classes.timelineItem}>
                <div className={`${classes.timelineDot} ${classes.name}`}></div>
                <span className={classes.timelineLabel}>Full Name</span>
                <span className={classes.timelineValue}>{user.name}</span>
              </div>
              <div className={classes.timelineItem}>
                <div className={`${classes.timelineDot} ${classes.email}`}></div>
                <span className={classes.timelineLabel}>Email Address</span>
                <span className={classes.timelineValue}>{user.email}</span>
              </div>
              <div className={classes.timelineItem}>
                <div className={`${classes.timelineDot} ${classes.pass}`}></div>
                <span className={classes.timelineLabel}>Password</span>
                <span className={classes.timelineValue}>********</span>
              </div>
              <div className={classes.timelineItem}>
                <div className={`${classes.timelineDot} ${classes.tz}`}></div>
                <span className={classes.timelineLabel}>Timezone</span>
                <span className={classes.timelineValue}>{settings.timezone}</span>
              </div>
            </div>

            <div className={classes.actions}>
              <button className={`${classes.btn} ${classes.btnPrimary}`} onClick={handleEditClick}>
                Edit Profile
              </button>
              <button className={`${classes.btn} ${classes.btnDanger}`} onClick={handleDeleteUser}>
                Delete Account
              </button>
            </div>

            <button className={`${classes.btn} ${classes.btnBack}`} onClick={() => navigate('/')}>
              ← Back to Tasks
            </button>
          </>
        ) : (
          <form className={classes.editForm} onSubmit={handleUpdateUser}>
            <div className={classes.formGroup}>
              <label className={classes.formLabel}>Name</label>
              <input type="text" name="name" value={editForm.name}
                onChange={e => setEditForm(prev => ({ ...prev, name: e.target.value }))}
                className={classes.formInput} required />
            </div>
            <div className={classes.formGroup}>
              <label className={classes.formLabel}>Email</label>
              <input type="email" name="email" value={editForm.email}
                onChange={e => setEditForm(prev => ({ ...prev, email: e.target.value }))}
                className={classes.formInput} required />
            </div>
            <div className={classes.formGroup}>
              <label className={classes.formLabel}>New Password</label>
              <input type="password" name="password" value={editForm.password}
                onChange={e => setEditForm(prev => ({ ...prev, password: e.target.value }))}
                className={classes.formInput} placeholder="Leave blank to keep current" />
            </div>
            <div className={classes.formGroup}>
              <label className={classes.formLabel}>Timezone</label>
              <select name="timezone" value={editForm.timezone}
                onChange={e => setEditForm(prev => ({ ...prev, timezone: e.target.value }))}
                className={classes.formSelect}>
                {timezoneOptions.map(tz => <option key={tz} value={tz}>{tz}</option>)}
              </select>
            </div>
            <div className={classes.formActions}>
              <button type="submit" className={`${classes.btn} ${classes.btnPrimary}`}>Save Changes</button>
              <button type="button" className={`${classes.btn} ${classes.btnDanger}`} onClick={handleCancelEdit}>Cancel</button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
};

export default ProfilePage;