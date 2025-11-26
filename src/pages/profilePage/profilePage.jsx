import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import classes from "./profile.module.css";
import LoadingSpinner from "../../components/loadingSpinner/loadingSpinner"

const ProfilePage = () => {
  const navigate = useNavigate();
   
  const [user, setUser] = useState({
    id: '',
    name: '',
    email: '',
    password: ''
  });

  const [settings, setSettings] = useState({
    timezone: 'UTC',
    duration: 60
  });

  const [isEditing, setIsEditing] = useState(false);
  const [editForm, setEditForm] = useState({
    name: '',
    email: '',
    password: '',
    timezone: 'UTC',
    duration: ''
  });

  const [originalData, setOriginalData] = useState({
    user: {
      name: '',
      email: '',
      password: ''
    },
    settings: {
      timezone: 'UTC',
      duration: 60
    }
  });

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // Список доступных временных зон
  const timezoneOptions = [
    'UTC-12', 'UTC-11', 'UTC-10', 'UTC-9', 'UTC-8', 'UTC-7', 'UTC-6', 'UTC-5', 
    'UTC-4', 'UTC-3', 'UTC-2', 'UTC-1', 'UTC', 'UTC+1', 'UTC+2', 'UTC+3', 
    'UTC+4', 'UTC+5', 'UTC+6', 'UTC+7', 'UTC+8', 'UTC+9', 'UTC+10', 'UTC+11', 'UTC+12'
  ];

  // Функция для конвертации минут в формат "дни, часы:минуты"
  const formatMinutesToTime = (minutes) => {
    const days = Math.floor(minutes / (24 * 60));
    const hours = Math.floor((minutes % (24 * 60)) / 60);
    const mins = minutes % 60;
    
    if (days > 0) {
      return `${days}d ${hours}h ${mins.toString().padStart(2, '0')}m`;
    } else {
      return `${hours}h ${mins.toString().padStart(2, '0')}m`;
    }
  };

  // Функция для конвертации формата "дни, часы:минуты" в минуты
  const formatTimeToMinutes = (timeString) => {
    let totalMinutes = 0;
    
    // Обработка дней
    const daysMatch = timeString.match(/(\d+)d/);
    if (daysMatch) {
      totalMinutes += parseInt(daysMatch[1]) * 24 * 60;
    }
    
    // Обработка часов
    const hoursMatch = timeString.match(/(\d+)h/);
    if (hoursMatch) {
      totalMinutes += parseInt(hoursMatch[1]) * 60;
    }
    
    // Обработка минут
    const minutesMatch = timeString.match(/(\d+)m/);
    if (minutesMatch) {
      totalMinutes += parseInt(minutesMatch[1]);
    }
    
    return totalMinutes;
  };

  // Загрузка данных пользователя и настроек
  useEffect(() => {
    const fetchUserData = async () => {
      try {
        setLoading(true);
        
        // Запрос данных пользователя
        const userResponse = await fetch('http://localhost:8080/v1/private/users', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: "include"
        });

        if (!userResponse.ok) {
          throw new Error('Ошибка при загрузке данных пользователя');
        }

        const userData = await userResponse.json();
        console.log('User data:', userData);
        
        // Запрос настроек
        const settingsResponse = await fetch('http://localhost:8080/v1/private/settings', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: "include"
        });

        if (!settingsResponse.ok) {
          throw new Error('Ошибка при загрузке настроек');
        }

        const settingsData = await settingsResponse.json();
        console.log('Settings data:', settingsData);

        const originalUserData = {
          name: userData.user.user_name || '',
          email: userData.user.user_email || '',
          password: ''
        };

        const originalSettingsData = {
          timezone: settingsData.settings.default_tz || 'UTC',
          duration: settingsData.settings.default_duration || 60
        };

        setOriginalData({
          user: originalUserData,
          settings: originalSettingsData
        });

        setUser({
          id: userData.id || '',
          name: userData.user.user_name || '',
          email: userData.user.user_email || '',
          password: '********'
        });

        setSettings({
          timezone: settingsData.settings.default_tz || 'UTC',
          duration: settingsData.settings.default_duration || 60
        });

        setEditForm({
          name: userData.user.user_name || '',
          email: userData.user.user_email || '',
          password: '',
          timezone: settingsData.settings.default_tz || 'UTC',
          duration: formatMinutesToTime(settingsData.settings.default_duration || 60)
        });

      } catch (error) {
        console.error('Ошибка:', error);
        setError(error.message);
      } finally {
        setLoading(false);
      }
    };

    fetchUserData();
  }, []);

  const handleGoHome = () => {
    navigate('/');
  };

  const handleDeleteUser = async () => {
    if (window.confirm('Вы уверены, что хотите удалить пользователя?')) {
      try {
        const response = await fetch(`http://localhost:8080/v1/private/users`, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: "include"
        });

        if (response.ok) {
          alert('Пользователь успешно удален');
          navigate('/auth');
        } else {
          alert('Ошибка при удалении пользователя');
        }
      } catch (error) {
        console.error('Ошибка:', error);
        alert('Произошла ошибка при удалении пользователя');
      }
    }
  };

  // Функция для подготовки данных пользователя к отправке
  const prepareUserData = () => {
    const userData = {
      user_name: editForm.name !== originalData.user.name ? editForm.name : null,
      user_email: editForm.email !== originalData.user.email ? editForm.email : null,
      user_pass: editForm.password ? editForm.password : null
    };

    if (editForm.password && Object.values(userData).every(val => val === null)) {
      userData.user_pass = editForm.password;
    }

    return userData;
  };

  const prepareSettingsData = () => {
    const durationInMinutes = formatTimeToMinutes(editForm.duration);
    
    return {
      default_tz: editForm.timezone !== originalData.settings.timezone ? editForm.timezone : null,
      default_duration: durationInMinutes !== originalData.settings.duration ? 
        durationInMinutes : null
    };
  };

  const handleUpdateUser = async (e) => {
    e.preventDefault();
    try {
      const userData = prepareUserData();
      const settingsData = prepareSettingsData();

      console.log('Sending user update:', userData);
      console.log('Sending settings update:', settingsData);

      let userResponse = null;
      let settingsResponse = null;

      const hasUserChanges = Object.values(userData).some(val => val !== null);
      if (hasUserChanges) {
        userResponse = await fetch(`http://localhost:8080/v1/private/users`, {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: "include",
          body: JSON.stringify(userData),
        });

        console.log('User response status:', userResponse?.status);
        
        if (!userResponse.ok) {
          const errorText = await userResponse.text();
          console.error('User update error:', errorText);
          throw new Error(`Ошибка при обновлении данных пользователя: ${userResponse.status}`);
        }
      }

      const hasSettingsChanges = Object.values(settingsData).some(val => val !== null);
      if (hasSettingsChanges) {
        settingsResponse = await fetch('http://localhost:8080/v1/private/settings', {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: "include",
          body: JSON.stringify(settingsData),
        });

        console.log('Settings response status:', settingsResponse?.status);

        if (!settingsResponse.ok) {
          const errorText = await settingsResponse.text();
          console.error('Settings update error:', errorText);
          throw new Error(`Ошибка при обновлении настроек: ${settingsResponse.status}`);
        }
      }

      if (!hasUserChanges && !hasSettingsChanges) {
        setIsEditing(false);
        alert('Нет изменений для сохранения');
        return;
      }

      if (hasUserChanges) {
        const updatedUser = userResponse ? await userResponse.json() : null;
        console.log('Updated user data:', updatedUser);

        setUser(prev => ({
          ...prev,
          name: editForm.name,
          email: editForm.email
        }));

        setOriginalData(prev => ({
          ...prev,
          user: {
            name: editForm.name,
            email: editForm.email,
            password: ''
          }
        }));
      }

      if (hasSettingsChanges) {
        const updatedSettings = settingsResponse ? await settingsResponse.json() : null;
        console.log('Updated settings data:', updatedSettings);

        const durationInMinutes = formatTimeToMinutes(editForm.duration);
        
        setSettings({
          timezone: editForm.timezone,
          duration: durationInMinutes
        });

        setOriginalData(prev => ({
          ...prev,
          settings: {
            timezone: editForm.timezone,
            duration: durationInMinutes
          }
        }));
      }

      setIsEditing(false);
      alert('Данные пользователя обновлены');

    } catch (error) {
      console.error('Ошибка:', error);
      alert(`Произошла ошибка при обновлении данных: ${error.message}`);
    }
  };

  const handleEditClick = () => {
    setEditForm({
      name: user.name,
      email: user.email,
      password: '',
      timezone: settings.timezone,
      duration: formatMinutesToTime(settings.duration)
    });
    setIsEditing(true);
  };

  const handleCancelEdit = () => {
    setIsEditing(false);
    setEditForm({
      name: user.name,
      email: user.email,
      password: '',
      timezone: settings.timezone,
      duration: formatMinutesToTime(settings.duration)
    });
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setEditForm(prev => ({
      ...prev,
      [name]: value
    }));
  };

  if (loading) {
    return (
      <div className={classes.loadingContainer}>
        <LoadingSpinner/>
        <div className={classes.loadingText}>Loading profile...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={classes.errorContainer}>
        <div className={classes.errorText}>Error: {error}</div>
        <button className={classes.retryButton} onClick={() => window.location.reload()}>
          Try again
        </button>
      </div>
    );
  }

  return (
    <div className={classes.userProfile}>
      <div className={classes.profileBackground}>
        <div className={classes.floatingShape1}></div>
        <div className={classes.floatingShape2}></div>
        <div className={classes.floatingShape3}></div>
      </div>
      
      <div className={classes.profileContainer}>
        <div className={classes.profileHeader}>
          <div className={classes.avatar}>
            {user.name ? user.name.charAt(0).toUpperCase() : 'U'}
          </div>
          <h2 className={classes.title}>My Profile</h2>
        </div>
        
        {!isEditing ? (
          // Режим просмотра
          <div className={classes.userInfo}>
            <div className={classes.infoSection}>
              <h3 className={classes.sectionTitle}>Personal Information</h3>
              <div className={classes.infoGrid}>
                <div className={classes.infoCard}>
                  <div className={classes.infoContent}>
                    <label className={classes.label}>Name</label>
                    <span className={classes.value}>{user.name}</span>
                  </div>
                </div>
                
                <div className={classes.infoCard}>
                  <div className={classes.infoContent}>
                    <label className={classes.label}>Email</label>
                    <span className={classes.value}>{user.email}</span>
                  </div>
                </div>
                
                <div className={classes.infoCard}>
                  <div className={classes.infoContent}>
                    <label className={classes.label}>Password</label>
                    <span className={classes.passwordMasked}>{user.password}</span>
                  </div>
                </div>
              </div>
            </div>

            <div className={classes.infoSection}>
              <h3 className={classes.sectionTitle}>Settings</h3>
              <div className={classes.infoGrid}>
                <div className={classes.infoCard}>
                  <div className={classes.infoContent}>
                    <label className={classes.label}>Timezone</label>
                    <span className={classes.value}>{settings.timezone}</span>
                  </div>
                </div>
                
                <div className={classes.infoCard}>
                  <div className={classes.infoContent}>
                    <label className={classes.label}>Duration</label>
                    <span className={classes.value}>{formatMinutesToTime(settings.duration)}</span>
                  </div>
                </div>
              </div>
            </div>
            
            <div className={classes.buttonsContainer}>
              <button 
                className={`${classes.button} ${classes.editButton}`}
                onClick={handleEditClick}
              >
                Edit Profile
              </button>
              <button 
                className={`${classes.button} ${classes.deleteButton}`}
                onClick={handleDeleteUser}
              >
                Delete Me
              </button>
            </div>
            
            <div className={classes.homeButtonContainer}>
              <button 
                className={`${classes.button} ${classes.homeButton}`}
                onClick={handleGoHome}
              >
                To My Tasks
              </button>
            </div>
          </div>
        ) : (
          // Режим редактирования
          <form className={classes.editForm} onSubmit={handleUpdateUser}>
            <div className={classes.formSection}>
              <h3 className={classes.sectionTitle}>Edit Personal Information</h3>
              <div className={classes.formGrid}>
                <div className={classes.formGroup}>
                  <label htmlFor="name" className={classes.formLabel}>
                    Name
                  </label>
                  <input
                    type="text"
                    id="name"
                    name="name"
                    value={editForm.name}
                    onChange={handleInputChange}
                    className={classes.formInput}
                    required
                  />
                </div>
                
                <div className={classes.formGroup}>
                  <label htmlFor="email" className={classes.formLabel}>
                    Email
                  </label>
                  <input
                    type="email"
                    id="email"
                    name="email"
                    value={editForm.email}
                    onChange={handleInputChange}
                    className={classes.formInput}
                    required
                  />
                </div>
                
                <div className={classes.formGroup}>
                  <label htmlFor="password" className={classes.formLabel}>
                    New Password
                  </label>
                  <input
                    type="password"
                    id="password"
                    name="password"
                    value={editForm.password}
                    onChange={handleInputChange}
                    className={classes.formInput}
                    placeholder="Leave blank to keep current password"
                  />
                </div>
              </div>
            </div>

            <div className={classes.formSection}>
              <h3 className={classes.sectionTitle}>Edit Preferences</h3>
              <div className={classes.formGrid}>
                <div className={classes.formGroup}>
                  <label htmlFor="timezone" className={classes.formLabel}>
                    Timezone
                  </label>
                  <select
                    id="timezone"
                    name="timezone"
                    value={editForm.timezone}
                    onChange={handleInputChange}
                    className={classes.formSelect}
                  >
                    {timezoneOptions.map(tz => (
                      <option key={tz} value={tz}>{tz}</option>
                    ))}
                  </select>
                </div>

                <div className={classes.formGroup}>
                  <label htmlFor="duration" className={classes.formLabel}>
                    Duration (format: Xd Yh Zm)
                  </label>
                  <input
                    type="text"
                    id="duration"
                    name="duration"
                    value={editForm.duration}
                    onChange={handleInputChange}
                    className={classes.formInput}
                    placeholder="e.g., 1d 2h 30m or 2h 30m or 45m"
                    pattern="(\d+d\s*)?(\d+h\s*)?(\d+m)?"
                    title="Please enter duration in format: Xd Yh Zm (e.g., 1d 2h 30m, 2h 30m, or 45m)"
                    required
                  />
                </div>
              </div>
            </div>
            
            <div className={classes.formButtons}>
              <button type="submit" className={`${classes.button} ${classes.saveButton}`}>
                Save Changes
              </button>
              <button 
                type="button" 
                className={`${classes.button} ${classes.cancelButton}`}
                onClick={handleCancelEdit}
              >
                Cancel
              </button>
            </div>
            
            <div className={classes.homeButtonContainer}>
              <button 
                className={`${classes.button} ${classes.homeButton}`}
                onClick={handleGoHome}
              >
                To My Tasks
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
};

export default ProfilePage;