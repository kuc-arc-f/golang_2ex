import React, { useState, useEffect } from 'react';
import axios from 'axios';
import Modal from 'react-modal'; 

// Modalのスタイル
const customStyles = {
  content: {
    top: '50%',
    left: '50%',
    right: 'auto',
    bottom: 'auto',
    marginRight: '-50%',
    transform: 'translate(-50%, -50%)',
    width: '80%',
    maxWidth: '600px',
  },
};

Modal.setAppElement('#root');

interface Todo {
  id: number;
  title: string;
  content: string;
  completed: boolean;
  content_type: string;
  is_public: boolean;
  food_orange: boolean;
  food_apple: boolean;
  food_banana: boolean;
  food_melon: boolean;
  food_grape: boolean;
}

const App: React.FC = () => {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [modalIsOpen, setModalIsOpen] = useState(false);
  const [selectedTodo, setSelectedTodo] = useState<Todo | null>(null);

  useEffect(() => {
    fetchTodos();
  }, []);

  const fetchTodos = async () => {
    const response = await axios.get('http://localhost:8080/api/list');
    setTodos(response.data);
  };

  const openModal = (todo: Todo | null) => {
    setSelectedTodo(todo);
    setModalIsOpen(true);
  };

  const closeModal = () => {
    setSelectedTodo(null);
    setModalIsOpen(false);
  };

  const handleSave = async (todo: Todo) => {
    if (todo.id) {
      await axios.post('http://localhost:8080/api/update', todo);
    } else {
      await axios.post('http://localhost:8080/api/create', todo);
    }
    fetchTodos();
    closeModal();
  };

  const handleDelete = async (id: number) => {
    await axios.post('http://localhost:8080/api/delete', { id });
    fetchTodos();
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Todo List</h1>
      <button
        onClick={() => openModal({
          id: 0,
          title: '',
          content: '',
          completed: false,
          content_type: '',
          is_public: false,
          food_orange: false,
          food_apple: false,
          food_banana: false,
          food_melon: false,
          food_grape: false,
        })}
        className="bg-blue-500 text-white px-4 py-2 rounded mb-4"
      >
        Add Todo
      </button>
      <table className="min-w-full bg-white">
        <thead>
          <tr>
            <th className="py-2">Title</th>
            <th className="py-2">Content</th>
            <th className="py-2">Completed</th>
            <th className="py-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          {todos.map((todo) => (
            <tr key={todo.id}>
              <td className="border px-4 py-2">{todo.title}</td>
              <td className="border px-4 py-2">{todo.content}</td>
              <td className="border px-4 py-2">{todo.completed ? 'Yes' : 'No'}</td>
              <td className="border px-4 py-2">
                <button
                  onClick={() => openModal(todo)}
                  className="bg-green-500 text-white px-2 py-1 rounded mr-2"
                >
                  Edit
                </button>
                <button
                  onClick={() => handleDelete(todo.id)}
                  className="bg-red-500 text-white px-2 py-1 rounded"
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {modalIsOpen && (
        <Modal
          isOpen={modalIsOpen}
          onRequestClose={closeModal}
          style={customStyles}
          contentLabel="Todo Modal"
        >
          <TodoForm
            todo={selectedTodo}
            onSave={handleSave}
            onCancel={closeModal}
          />
        </Modal>
      )}
    </div>
  );
};

interface TodoFormProps {
  todo: Todo | null;
  onSave: (todo: Todo) => void;
  onCancel: () => void;
}

const TodoForm: React.FC<TodoFormProps> = ({ todo, onSave, onCancel }) => {
  const [formState, setFormState] = useState<Todo>(
    todo || {
      id: 0,
      title: '',
      content: '',
      completed: false,
      content_type: '',
      is_public: false,
      food_orange: false,
      food_apple: false,
      food_banana: false,
      food_melon: false,
      food_grape: false,
    }
  );

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target;
    setFormState((prevState) => ({
      ...prevState,
      [name]: type === 'checkbox' ? checked : value,
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave(formState);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2">
          Title
        </label>
        <input
          type="text"
          name="title"
          value={formState.title}
          onChange={handleChange}
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
        />
      </div>
      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2">
          Content
        </label>
        <input
          type="text"
          name="content"
          value={formState.content}
          onChange={handleChange}
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
        />
      </div>
      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2">
          Content Type
        </label>
        <input
          type="text"
          name="content_type"
          value={formState.content_type}
          onChange={handleChange}
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
        />
      </div>
      <div className="mb-4">
        <label className="inline-flex items-center">
          <input
            type="checkbox"
            name="completed"
            checked={formState.completed}
            onChange={handleChange}
            className="form-checkbox"
          />
          <span className="ml-2">Completed</span>
        </label>
      </div>
      <div className="mb-4">
        <span className="block text-gray-700 text-sm font-bold mb-2">
          Public/Private
        </span>
        <label className="inline-flex items-center">
          <input
            type="radio"
            name="is_public"
            value="true"
            checked={formState.is_public === true}
            onChange={() =>
              setFormState((prevState) => ({ ...prevState, is_public: true }))
            }
            className="form-radio"
          />
          <span className="ml-2">Public</span>
        </label>
        <label className="inline-flex items-center ml-6">
          <input
            type="radio"
            name="is_public"
            value="false"
            checked={formState.is_public === false}
            onChange={() =>
              setFormState((prevState) => ({ ...prevState, is_public: false }))
            }
            className="form-radio"
          />
          <span className="ml-2">Private</span>
        </label>
      </div>
      <div className="mb-4">
        <span className="block text-gray-700 text-sm font-bold mb-2">
          Favorite Foods
        </span>
        <label className="inline-flex items-center">
          <input
            type="checkbox"
            name="food_orange"
            checked={formState.food_orange}
            onChange={handleChange}
            className="form-checkbox"
          />
          <span className="ml-2">Orange</span>
        </label>
        <label className="inline-flex items-center ml-6">
          <input
            type="checkbox"
            name="food_apple"
            checked={formState.food_apple}
            onChange={handleChange}
            className="form-checkbox"
          />
          <span className="ml-2">Apple</span>
        </label>
        <label className="inline-flex items-center ml-6">
          <input
            type="checkbox"
            name="food_banana"
            checked={formState.food_banana}
            onChange={handleChange}
            className="form-checkbox"
          />
          <span className="ml-2">Banana</span>
        </label>
        <label className="inline-flex items-center ml-6">
          <input
            type="checkbox"
            name="food_melon"
            checked={formState.food_melon}
            onChange={handleChange}
            className="form-checkbox"
          />
          <span className="ml-2">Melon</span>
        </label>
        <label className="inline-flex items-center ml-6">
          <input
            type="checkbox"
            name="food_grape"
            checked={formState.food_grape}
            onChange={handleChange}
            className="form-checkbox"
          />
          <span className="ml-2">Grape</span>
        </label>
      </div>
      <div className="flex items-center justify-between">
        <button
          type="submit"
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
        >
          Save
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
        >
          Cancel
        </button>
      </div>
    </form>
  );
};

export default App;