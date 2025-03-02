import React, { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from '../AxiosInstance';
import './EmployeeDetail.css';

function EmployeeDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [employee, setEmployee] = useState({ name: '', position: '', contact: '' });

  const fetchEmployee = useCallback(async () => {
    try {
      const response = await axios.get('https://localhost:5000/crud/FetchTheDetails');
      const emp = response.data.find((e) => e._id === id);
      if (emp) {
        setEmployee(emp);
      } else {
        console.warn('Employee not found');
      }
    } catch (error) {
      console.error('Failed to fetch employee details:', error);
    }
  }, [id]);

  useEffect(() => {
    fetchEmployee();
  }, [fetchEmployee]);

  const handleUpdate = async () => {
    try {
      const updateFields = { name: employee.name, position: employee.position, contact: Number(employee.contact) };
      await axios.put(`http://localhosts:5000/crud/UpdateTheDetails/${id}`, updateFields);
      navigate('/home');
    } catch (error) {
      console.error('Failed to update employee:', error);
    }
  };

  const handleDelete = async () => {
    try {
      await axios.delete(`http://localhosts:5000/crud/DeleteTheParticularDetails/${id}`);
      navigate('/home');
    } catch (error) {
      console.error('Failed to delete employee:', error);
    }
  };

  return (
    <div className="detail-container">
      <div className="detail-box">
        <h2>Edit Employee Details</h2>
        <div className="form-group">
          <label>Name</label>
          <input
            type="text"
            value={employee.name}
            onChange={(e) => setEmployee({ ...employee, name: e.target.value })}
          />
        </div>
        <div className="form-group">
          <label>Position</label>
          <input
            type="text"
            value={employee.position}
            onChange={(e) => setEmployee({ ...employee, position: e.target.value })}
          />
        </div>
        <div className="form-group">
          <label>Contact</label>
          <input
            type="text"
            value={employee.contact}
            onChange={(e) => setEmployee({ ...employee, contact: e.target.value })}
          />
        </div>
        <div className="button-group">
          <button onClick={handleUpdate} className="update">Update</button>
          <button onClick={handleDelete} className="delete">Delete</button>
        </div>
      </div>
    </div>
  );
}

export default EmployeeDetail;
