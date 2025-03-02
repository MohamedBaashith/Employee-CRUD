import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from '../AxiosInstance';
import './Login.css';

function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await axios.post('https://localhost:5000/auth/Login', { username, password });
      console.log("Login Response:", response.data);
      if (response.data.success === true || response.data.success === "true") {
        localStorage.setItem('token', response.data.token);
        navigate('/home');
      } else {
        setError('Invalid Credentials');
      }
    } catch (err) {
      setError('An error occurred');
    }
  };
  return (
    <div className="container">
      <div className="login-box">
        <h2>Login</h2>
        {error && <p className="error">{error}</p>}
        <form onSubmit={handleSubmit}>
          <div className="input-group">
            <label>Username</label>
            <input
              type="text"
              placeholder="Enter Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </div>
          <div className="input-group">
            <label>Password</label>
            <input
              type="password"
              placeholder="Enter Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>
          <div>
            <Link to="/register">New Employee??? Register!!!</Link>
          </div>
          <div className="button-group">
            <button type="submit">Login</button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default Login;
