import React from 'react';
import { TodoItem } from '../services/api';
import { CheckCircle2, Clock, AlertCircle, User, Tag } from 'lucide-react';

interface TodoCardProps {
  todo: TodoItem;
}

const TodoCard: React.FC<TodoCardProps> = ({ todo }) => {
  const getPriorityColor = (priority: string) => {
    switch (priority.toLowerCase()) {
      case 'high':
        return 'text-red-500 border-red-500/30 bg-red-500/10';
      case 'medium':
        return 'text-yellow-500 border-yellow-500/30 bg-yellow-500/10';
      case 'low':
        return 'text-blue-500 border-blue-500/30 bg-blue-500/10';
      default:
        return 'text-gray-500 border-gray-500/30 bg-gray-500/10';
    }
  };

  const getPriorityIcon = (priority: string) => {
    switch (priority.toLowerCase()) {
      case 'high':
        return <AlertCircle size={16} className="mr-1" />;
      case 'medium':
        return <Clock size={16} className="mr-1" />;
      case 'low':
        return <CheckCircle2 size={16} className="mr-1" />;
      default:
        return null;
    }
  };

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-lg p-4 hover:border-gray-700 transition-colors duration-200">
      <div className="flex justify-between items-start mb-2">
        <h3 className="text-white font-semibold text-lg leading-tight">{todo.title}</h3>
        <span className={`flex items-center px-2 py-0.5 rounded text-xs font-medium border ${getPriorityColor(todo.priority)}`}>
          {getPriorityIcon(todo.priority)}
          {todo.priority.toUpperCase()}
        </span>
      </div>
      
      <p className="text-gray-400 text-sm mb-3 line-clamp-2">{todo.description}</p>
      
      <div className="flex flex-wrap gap-2 mb-3">
        {todo.tags && todo.tags.map((tag, index) => (
          <span key={index} className="flex items-center px-2 py-0.5 rounded-full bg-gray-800 text-gray-300 text-xs">
            <Tag size={10} className="mr-1" />
            {tag}
          </span>
        ))}
      </div>

      <div className="flex justify-between items-center text-xs text-gray-500 border-t border-gray-800 pt-3 mt-auto">
        <div className="flex items-center">
           {todo.deadline && (
            <span className="flex items-center mr-3 text-orange-400/80">
              <Clock size={12} className="mr-1" />
              {todo.deadline}
            </span>
          )}
          {todo.estimated_duration && (
             <span className="text-gray-600 mr-3">
               ⏱ {todo.estimated_duration}
             </span>
          )}
        </div>
        
        {todo.related_people && todo.related_people.length > 0 && (
          <div className="flex items-center">
            <User size={12} className="mr-1" />
            <span className="truncate max-w-[100px]">{todo.related_people.join(', ')}</span>
          </div>
        )}
      </div>
    </div>
  );
};

export default TodoCard;
