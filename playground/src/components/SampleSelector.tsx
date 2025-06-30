'use client';

import { useState } from 'react';
import { 
  qasmSamples, 
  sampleCategories, 
  difficultyColors, 
  difficultyLabels,
  type QASMSample 
} from '@/data/samples';

interface SampleSelectorProps {
  onSelectSample: (sample: QASMSample) => void;
  onClose: () => void;
  isOpen: boolean;
}

export default function SampleSelector({ onSelectSample, onClose, isOpen }: SampleSelectorProps) {
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [searchTerm, setSearchTerm] = useState('');

  const filteredSamples = qasmSamples.filter(sample => {
    const matchesCategory = selectedCategory === 'all' || sample.category === selectedCategory;
    const matchesSearch = sample.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         sample.description.toLowerCase().includes(searchTerm.toLowerCase());
    return matchesCategory && matchesSearch;
  });

  const handleSelectSample = (sample: QASMSample) => {
    onSelectSample(sample);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="modal modal-open">
      <div className="modal-box w-11/12 max-w-4xl h-[85vh] md:h-[80vh] flex flex-col">
        <div className="flex justify-between items-center mb-3 md:mb-4">
          <h3 className="font-bold text-base md:text-lg">Choose a Sample</h3>
          <button className="btn btn-sm btn-circle btn-ghost" onClick={onClose}>✕</button>
        </div>

        {/* Search and Filter */}
        <div className="flex flex-col gap-3 mb-3 md:mb-4">
          <div className="w-full">
            <input
              type="text"
              placeholder="Search samples..."
              className="input input-bordered input-sm md:input-md w-full"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </div>
          
          <div className="flex gap-1 md:gap-2 overflow-x-auto pb-1">
            {sampleCategories.map(category => (
              <button
                key={category.id}
                className={`btn btn-xs md:btn-sm whitespace-nowrap flex-shrink-0 ${
                  selectedCategory === category.id ? 'btn-primary' : 'btn-outline'
                }`}
                onClick={() => setSelectedCategory(category.id)}
              >
                <span className="text-xs md:text-sm">{category.icon}</span>
                <span className="hidden sm:inline text-xs md:text-sm">{category.name}</span>
              </button>
            ))}
          </div>
        </div>

        {/* Sample Grid */}
        <div className="flex-1 overflow-y-auto">
          {filteredSamples.length === 0 ? (
            <div className="text-center py-6 md:py-8">
              <div className="text-3xl md:text-4xl mb-2">🔍</div>
              <p className="text-base md:text-lg font-semibold">No samples found</p>
              <p className="text-xs md:text-sm opacity-70">Try adjusting your search or filters</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3 md:gap-4">
              {filteredSamples.map(sample => (
                <div
                  key={sample.id}
                  className="card bg-base-200 shadow-lg hover:shadow-xl transition-all duration-200 cursor-pointer hover:scale-[1.02]"
                  onClick={() => handleSelectSample(sample)}
                >
                  <div className="card-body p-3 md:p-4">
                    <div className="flex justify-between items-start mb-2">
                      <h4 className="card-title text-xs md:text-sm leading-tight flex-1 mr-2">{sample.name}</h4>
                      <div className={`badge badge-xs md:badge-sm flex-shrink-0 ${difficultyColors[sample.difficulty]}`}>
                        <span className="text-xs">{difficultyLabels[sample.difficulty]}</span>
                      </div>
                    </div>
                    
                    <p className="text-xs opacity-80 mb-2 md:mb-3 line-clamp-2 leading-relaxed">
                      {sample.description}
                    </p>
                    
                    <div className="card-actions justify-between items-center">
                      <div className="text-xs opacity-60 flex items-center">
                        <span className="mr-1">{sampleCategories.find(cat => cat.id === sample.category)?.icon}</span>
                        <span className="hidden sm:inline text-xs">
                          {sampleCategories.find(cat => cat.id === sample.category)?.name}
                        </span>
                      </div>
                      
                      <button className="btn btn-primary btn-xs md:btn-sm">
                        Load
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="modal-action mt-3 md:mt-4">
          <button className="btn btn-outline btn-sm md:btn-md" onClick={onClose}>Cancel</button>
        </div>
      </div>
    </div>
  );
}