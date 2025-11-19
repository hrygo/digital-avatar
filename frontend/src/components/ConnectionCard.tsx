import React from 'react';
import { ConnectionAnalysis } from '../services/api';
import { Users, MessageSquare, TrendingUp, AlertTriangle, Zap } from 'lucide-react';

interface ConnectionCardProps {
  connection: ConnectionAnalysis;
}

const ConnectionCard: React.FC<ConnectionCardProps> = ({ connection }) => {
  const getImportanceColor = (importance: string) => {
    switch (importance?.toLowerCase()) {
      case 'high':
        return 'text-red-400 bg-red-400/10 border-red-400/20';
      case 'medium':
        return 'text-yellow-400 bg-yellow-400/10 border-yellow-400/20';
      case 'low':
        return 'text-blue-400 bg-blue-400/10 border-blue-400/20';
      default:
        return 'text-gray-400 bg-gray-400/10 border-gray-400/20';
    }
  };

  const getSentimentColor = (sentiment: string) => {
    switch (sentiment?.toLowerCase()) {
      case 'positive':
        return 'text-green-500';
      case 'negative':
        return 'text-red-500';
      default:
        return 'text-gray-500';
    }
  };

  return (
    <div className="bg-gray-900/50 border border-gray-800 rounded-xl p-5 hover:bg-gray-800/50 hover:border-gray-700 transition-all duration-300 group">
      <div className="flex items-start justify-between mb-4">
        <div className="flex items-center space-x-3">
          <div className="w-10 h-10 rounded-full bg-gradient-to-br from-purple-500 to-indigo-600 flex items-center justify-center text-white font-bold text-lg shadow-lg">
            {connection.person.charAt(0)}
          </div>
          <div>
            <h3 className="text-white font-semibold text-base">{connection.person}</h3>
            <div className="flex items-center space-x-2 text-xs mt-0.5">
              <span className="text-gray-400">{connection.relationship_type}</span>
              <span className="text-gray-600">•</span>
              <span className="text-gray-400">{connection.interaction_pattern}</span>
            </div>
          </div>
        </div>
        <span className={`px-2 py-1 rounded text-xs font-medium border ${getImportanceColor(connection.importance)}`}>
          {connection.importance?.toUpperCase() || 'NORMAL'}
        </span>
      </div>

      <div className="space-y-3">
        <div className="bg-gray-800/50 rounded-lg p-3 text-sm text-gray-300 border border-gray-700/50">
          <div className="flex items-center text-xs text-gray-500 mb-1">
            <MessageSquare size={12} className="mr-1" />
            <span>关键互动: {connection.action}</span>
          </div>
          <p className="line-clamp-2 leading-relaxed">{connection.context}</p>
        </div>

        <div className="grid grid-cols-2 gap-2 text-xs">
            <div className="flex items-center space-x-1.5 text-gray-400 bg-gray-800/30 p-2 rounded">
                <TrendingUp size={12} />
                <span>机会: {(connection.opportunity_score * 100).toFixed(0)}%</span>
            </div>
            <div className="flex items-center space-x-1.5 text-gray-400 bg-gray-800/30 p-2 rounded">
                <Zap size={12} className={getSentimentColor(connection.sentiment)} />
                <span>情感: {connection.sentiment}</span>
            </div>
        </div>

        {connection.suggested_action && (
          <div className="flex items-start space-x-2 text-xs text-blue-400 bg-blue-500/5 p-2 rounded border border-blue-500/10">
            <AlertTriangle size={12} className="mt-0.5 flex-shrink-0" />
            <span>建议: {connection.suggested_action}</span>
          </div>
        )}
        
         <div className="flex flex-wrap gap-1.5 pt-1">
            {connection.topics && connection.topics.map((topic, i) => (
                <span key={i} className="px-1.5 py-0.5 rounded bg-gray-800 text-gray-500 text-[10px]">
                    #{topic}
                </span>
            ))}
        </div>
      </div>
      
      <div className="mt-4 pt-3 border-t border-gray-800 flex justify-between items-center text-xs text-gray-500">
         <span>{connection.message_count} 次互动</span>
         <span>上次: {connection.last_interaction || '近期'}</span>
      </div>
    </div>
  );
};

export default ConnectionCard;
