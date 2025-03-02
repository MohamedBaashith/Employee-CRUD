import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from '../AxiosInstance';
import './Login.css';

function Register() {
  const [FirstName, setFirstName] = useState('');
  const [LastName, setLastName] = useState('');
  const [Email, setEmail] = useState('');
  const [Password, setNewPassword] = useState('');
  const [ConfirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (Password !== ConfirmPassword) {
      setError('Passwords do not match');
      return;
    }

    try {
      const response = await axios.post('https://localhost:5000/auth/Register', {
        FirstName,
        LastName,
        Email,
        Password,
      });
      console.log("Registration Response:", response.data);
      if (response.data.success === true || response.data.success === "true") {
        alert('Employee registered successfully!');
        navigate('/');
      } else {
        setError('Registration failed. Please try again.');
      }
    } catch (err) {
      alert("Email already exists!!!")
    }
  };

  return (
    <div className="container">
      <div className="login-box">
        <h2>Register</h2>
        {error && <p className="error">{error}</p>}
        <form onSubmit={handleSubmit}>
          <div className="input-group">
            <label>First Name</label>
            <input
              type="text"
              placeholder="Enter First Name"
              value={FirstName}
              onChange={(e) => setFirstName(e.target.value)}
            />
          </div>
          <div className="input-group">
            <label>Last Name</label>
            <input
              type="text"
              placeholder="Enter Last Name"
              value={LastName}
              onChange={(e) => setLastName(e.target.value)}
            />
          </div>
          <div className="input-group">
            <label>Email</label>
            <input
              type="text"
              placeholder="Enter Email"
              value={Email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </div>
          <div className="input-group">
            <label>New Password</label>
            <input
              type="password"
              placeholder="Enter New Password"
              value={Password}
              onChange={(e) => setNewPassword(e.target.value)}
            />
          </div>
          <div className="input-group">
            <label>Confirm Password</label>
            <input
              type="password"
              placeholder="Enter Confirm Password"
              value={ConfirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
            />
          </div>
          <div className="button-group">
            <button type="submit">Register</button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default Register;
