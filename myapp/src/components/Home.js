import React, { useState, useEffect } from 'react';
import axios from '../AxiosInstance';
import { Link } from 'react-router-dom';
import './Home.css';

function Home() {
  const [employees, setEmployees] = useState([]);
  const [name, setName] = useState('');
  const [position, setPosition] = useState('');
  const [contact, setContact] = useState('');

  useEffect(() => {
    fetchEmployees();
  }, []);

  const fetchEmployees = async () => {
    try {
      const response = await axios.get('https://localhost:5000/crud/FetchTheDetails');
      
      setEmployees(Array.isArray(response.data) ? response.data : []);
    } catch (error) {
      console.error('Failed to fetch employees:', error);
      setEmployees([]);
    }
  };

  const handleAddEmployee = async () => {
    const newEmployee = { name, position, contact: Number(contact) };
    try {
      await axios.post('https://localhost:5000/crud/AddTheDetails', newEmployee);
      fetchEmployees();
      setName('');
      setPosition('');
      setContact('');
    } catch (error) {
      console.error('Failed to add employee:', error);
    }
  };

  return (
    <div className="home-container">

      <div className="add-employee-form">
        <h2>Add New Employee Details</h2>
        <div className="form-group">
          <label>Name</label>
          <input
            type="text"
            placeholder="Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>
        <div className="form-group">
          <label>Position</label>
          <input
            type="text"
            placeholder="Position"
            value={position}
            onChange={(e) => setPosition(e.target.value)}
          />
        </div>
        <div className="form-group">
          <label>Contact</label>
          <input
            type="text"
            placeholder="Contact"
            value={contact}
            onChange={(e) => setContact(e.target.value)}
          />
        </div>
        <button onClick={handleAddEmployee}>Add Employee</button>
      </div>

      <div className="header">
        <h1>Employees List</h1>
      </div>
      <ul className="employee-list">
        {employees.map((employee) => (
          <li key={employee._id} className="employee-item">
            <div>
              <h3>{employee.name}</h3>
              <p>{employee.position} - {employee.contact}</p>
            </div>
            <div className="employee-actions">
              <Link to={`/employee/${employee._id}`}>Edit</Link>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default Home;
