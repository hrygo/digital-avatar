import React, { useState, useCallback } from 'react';
import {
  Database,
  FileText,
  Upload,
  ArrowRight,
  ArrowLeft,
  CheckCircle,
  AlertCircle,
  Loader2,
  Eye,
  EyeOff,
  FolderOpen,
  HardDrive
} from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { wechatService, ImportRequest, ImportResponse, WeChatStatus } from '../services/wechatService';

interface WeChatImportWizardProps {
  onComplete?: (response: ImportResponse) => void;
  onCancel?: () => void;
}

type ImportStep = 'source-select' | 'path-config' | 'importing' | 'completed' | 'error';
type ImportSource = 'database' | 'backup' | 'file';

const WeChatImportWizard: React.FC<WeChatImportWizardProps> = ({ onComplete, onCancel }) => {
  const [currentStep, setCurrentStep] = useState<ImportStep>('source-select');
  const [importSource, setImportSource] = useState<ImportSource>('database');
  const [importPath, setImportPath] = useState<string>('');
  const [importMode, setImportMode] = useState<'full' | 'incremental'>('full');
  const [importResponse, setImportResponse] = useState<ImportResponse | null>(null);
  const [errorMessage, setErrorMessage] = useState<string>('');
  const [wechatStatus, setWeChatStatus] = useState<WeChatStatus | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [showAdvancedOptions, setShowAdvancedOptions] = useState<boolean>(false);
  const [pathValidation, setPathValidation] = useState<{
    isValid: boolean;
    message: string;
  }>({ isValid: false, message: '' });

  const steps = [
    { id: 'source-select', title: '选择数据源', description: '选择微信数据的来源类型' },
    { id: 'path-config', title: '配置路径', description: '配置数据文件路径和导入选项' },
    { id: 'importing', title: '导入数据', description: '正在处理微信数据' },
    { id: 'completed', title: '导入完成', description: '数据导入成功完成' },
    { id: 'error', title: '导入失败', description: '数据导入过程中出现问题' }
  ];

  const sourceOptions = [
    {
      type: 'database' as ImportSource,
      title: '微信数据库',
      description: '直接连接微信数据库文件',
      icon: <Database className="w-6 h-6" />,
      recommended: true,
      features: ['实时数据访问', '最完整的聊天记录', '支持增量同步']
    },
    {
      type: 'backup' as ImportSource,
      title: '备份文件',
      description: '从微信备份文件中恢复',
      icon: <FileText className="w-6 h-6" />,
      recommended: false,
      features: ['支持微信备份格式', '安全的数据恢复', '支持多种备份版本']
    },
    {
      type: 'file' as ImportSource,
      title: '导出文件',
      description: '导入已导出的聊天记录',
      icon: <Upload className="w-6 h-6" />,
      recommended: false,
      features: ['支持多种导出格式', '灵活的文件选择', '自定义数据处理']
    }
  ];

  const validatePath = useCallback(async (path: string, source: ImportSource) => {
    if (!path.trim()) {
      setPathValidation({ isValid: false, message: '请输入数据路径' });
      return;
    }

    setIsLoading(true);
    try {
      // 这里可以添加路径验证逻辑
      // 暂时使用基本的路径格式验证
      const isValidFormat = path.length > 3 && (path.includes('/') || path.includes('\\'));

      if (isValidFormat) {
        setPathValidation({
          isValid: true,
          message: source === 'database' ? '数据库路径格式正确' : '文件路径格式正确'
        });
      } else {
        setPathValidation({
          isValid: false,
          message: '请输入有效的文件路径'
        });
      }
    } catch (error) {
      setPathValidation({
        isValid: false,
        message: '路径验证失败，请检查路径是否正确'
      });
    } finally {
      setIsLoading(false);
    }
  }, []);

  const handleSourceSelect = (source: ImportSource) => {
    setImportSource(source);
    setCurrentStep('path-config');
  };

  const handlePathConfirm = async () => {
    if (!pathValidation.isValid) {
      return;
    }

    setIsLoading(true);
    setCurrentStep('importing');

    try {
      const importRequest: ImportRequest = {
        type: importSource,
        path: importPath,
        mode: importMode
      };

      let response: ImportResponse;

      if (currentStep === 'path-config') {
        // 首次连接
        response = await wechatService.connect(importRequest);
      } else {
        // 同步操作
        response = await wechatService.sync(importRequest);
      }

      setImportResponse(response);

      // 获取微信状态
      try {
        const status = await wechatService.getStatus();
        setWeChatStatus(status);
      } catch (statusError) {
        console.warn('获取微信状态失败:', statusError);
      }

      setCurrentStep('completed');
      onComplete?.(response);
    } catch (error) {
      console.error('微信数据导入失败:', error);
      setErrorMessage(error instanceof Error ? error.message : '导入失败，请检查数据源配置');
      setCurrentStep('error');
    } finally {
      setIsLoading(false);
    }
  };

  const handleRetry = () => {
    setErrorMessage('');
    setImportResponse(null);
    setCurrentStep('source-select');
  };

  const handleBack = () => {
    if (currentStep === 'path-config') {
      setCurrentStep('source-select');
      setPathValidation({ isValid: false, message: '' });
    }
  };

  const handlePathChange = (value: string) => {
    setImportPath(value);
    if (pathValidation.isValid && !value.trim()) {
      setPathValidation({ isValid: false, message: '请输入数据路径' });
    }
  };

  const handlePathBlur = () => {
    validatePath(importPath, importSource);
  };

  const renderSourceSelect = () => (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      className="space-y-6"
    >
      <div className="text-center mb-8">
        <h2 className="text-2xl font-bold text-white mb-3">选择微信数据源</h2>
        <p className="text-gray-400">
          选择您要导入的微信数据类型，推荐使用数据库连接以获得最完整的数据
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {sourceOptions.map((option) => (
          <motion.button
            key={option.type}
            whileHover={{ scale: 1.02, y: -2 }}
            whileTap={{ scale: 0.98 }}
            onClick={() => handleSourceSelect(option.type)}
            className={`
              relative bg-gray-800/50 backdrop-blur-sm rounded-xl p-6 border-2 transition-all duration-200
              hover:border-blue-500/50 hover:bg-gray-700/50 text-left group
              ${importSource === option.type ? 'border-blue-500 bg-blue-600/10' : 'border-gray-700'}
            `}
          >
            {option.recommended && (
              <div className="absolute -top-2 -right-2 bg-blue-600 text-white text-xs px-2 py-1 rounded-full">
                推荐
              </div>
            )}

            <div className="flex items-center gap-3 mb-4 text-blue-400">
              {option.icon}
              <span className="font-semibold text-lg">{option.title}</span>
            </div>

            <p className="text-gray-400 text-sm mb-4">{option.description}</p>

            <div className="space-y-2">
              {option.features.map((feature, index) => (
                <div key={index} className="flex items-center gap-2 text-xs text-gray-500">
                  <CheckCircle className="w-3 h-3 text-green-400" />
                  <span>{feature}</span>
                </div>
              ))}
            </div>

            <div className="mt-4 flex items-center gap-2 text-blue-400 opacity-0 group-hover:opacity-100 transition-opacity">
              <span className="text-sm">选择此方案</span>
              <ArrowRight className="w-4 h-4" />
            </div>
          </motion.button>
        ))}
      </div>
    </motion.div>
  );

  const renderPathConfig = () => (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      className="space-y-6"
    >
      <div className="text-center mb-8">
        <h2 className="text-2xl font-bold text-white mb-3">配置数据路径</h2>
        <p className="text-gray-400">
          {importSource === 'database' ? '请输入微信数据库文件的完整路径' :
           importSource === 'backup' ? '请选择微信备份文件路径' :
           '请选择导出文件的路径'}
        </p>
      </div>

      <div className="bg-gray-800/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700">
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              数据路径
            </label>
            <div className="relative">
              <input
                type="text"
                value={importPath}
                onChange={(e) => handlePathChange(e.target.value)}
                onBlur={handlePathBlur}
                placeholder={
                  importSource === 'database'
                    ? '/Users/username/Library/Application Support/Tencent/MicroMsg/...'
                    : '/path/to/wechat/backup/or/export'
                }
                className={`
                  w-full px-4 py-3 pl-10 bg-gray-900/50 border rounded-lg text-white placeholder-gray-500
                  focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent
                  ${pathValidation.isValid ? 'border-green-500' : pathValidation.message ? 'border-red-500' : 'border-gray-600'}
                `}
              />
              <div className="absolute left-3 top-3.5 text-gray-400">
                {importSource === 'database' ? (
                  <HardDrive className="w-5 h-5" />
                ) : (
                  <FolderOpen className="w-5 h-5" />
                )}
              </div>
            </div>
            {pathValidation.message && (
              <div className={`mt-2 text-sm ${pathValidation.isValid ? 'text-green-400' : 'text-red-400'}`}>
                {pathValidation.message}
              </div>
            )}
          </div>

          <div>
            <button
              type="button"
              onClick={() => setShowAdvancedOptions(!showAdvancedOptions)}
              className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors"
            >
              {showAdvancedOptions ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
              <span className="text-sm">高级选项</span>
            </button>

            <AnimatePresence>
              {showAdvancedOptions && (
                <motion.div
                  initial={{ height: 0, opacity: 0 }}
                  animate={{ height: 'auto', opacity: 1 }}
                  exit={{ height: 0, opacity: 0 }}
                  transition={{ duration: 0.2 }}
                  className="overflow-hidden"
                >
                  <div className="mt-4 space-y-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-300 mb-2">
                        导入模式
                      </label>
                      <div className="grid grid-cols-2 gap-4">
                        <button
                          type="button"
                          onClick={() => setImportMode('full')}
                          className={`
                            px-4 py-2 rounded-lg border-2 transition-all
                            ${importMode === 'full'
                              ? 'border-blue-500 bg-blue-600/20 text-blue-400'
                              : 'border-gray-600 text-gray-400 hover:border-gray-500'}
                          `}
                        >
                          <div className="font-medium">全量导入</div>
                          <div className="text-xs opacity-70">导入所有历史数据</div>
                        </button>
                        <button
                          type="button"
                          onClick={() => setImportMode('incremental')}
                          className={`
                            px-4 py-2 rounded-lg border-2 transition-all
                            ${importMode === 'incremental'
                              ? 'border-blue-500 bg-blue-600/20 text-blue-400'
                              : 'border-gray-600 text-gray-400 hover:border-gray-500'}
                          `}
                        >
                          <div className="font-medium">增量导入</div>
                          <div className="text-xs opacity-70">仅导入新数据</div>
                        </button>
                      </div>
                    </div>
                  </div>
                </motion.div>
              )}
            </AnimatePresence>
          </div>
        </div>
      </div>

      <div className="flex justify-between">
        <button
          onClick={handleBack}
          className="flex items-center gap-2 px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
        >
          <ArrowLeft className="w-4 h-4" />
          返回
        </button>

        <button
          onClick={handlePathConfirm}
          disabled={!pathValidation.isValid || isLoading}
          className={`
            flex items-center gap-2 px-6 py-3 rounded-lg transition-colors
            ${pathValidation.isValid && !isLoading
              ? 'bg-blue-600 hover:bg-blue-700 text-white'
              : 'bg-gray-700 text-gray-400 cursor-not-allowed'}
          `}
        >
          {isLoading ? (
            <>
              <Loader2 className="w-4 h-4 animate-spin" />
              处理中...
            </>
          ) : (
            <>
              开始导入
              <ArrowRight className="w-4 h-4" />
            </>
          )}
        </button>
      </div>
    </motion.div>
  );

  const renderImporting = () => (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      className="text-center space-y-6 py-12"
    >
      <div className="flex justify-center">
        <Loader2 className="w-16 h-16 text-blue-500 animate-spin" />
      </div>

      <div>
        <h2 className="text-2xl font-bold text-white mb-3">正在导入微信数据</h2>
        <p className="text-gray-400">
          {importSource === 'database' ? '正在连接和解析微信数据库...' :
           importSource === 'backup' ? '正在恢复微信备份数据...' :
           '正在处理导出文件...'}
        </p>
      </div>

      <div className="bg-gray-800/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700 max-w-md mx-auto">
        <div className="space-y-3 text-left">
          <div className="flex items-center gap-3">
            <CheckCircle className="w-5 h-5 text-green-400" />
            <span className="text-gray-300">数据源验证完成</span>
          </div>
          <div className="flex items-center gap-3">
            <Loader2 className="w-5 h-5 text-blue-400 animate-spin" />
            <span className="text-gray-300">正在解析数据结构</span>
          </div>
          <div className="flex items-center gap-3">
            <div className="w-5 h-5 border-2 border-gray-600 rounded-full" />
            <span className="text-gray-500">提取消息记录</span>
          </div>
          <div className="flex items-center gap-3">
            <div className="w-5 h-5 border-2 border-gray-600 rounded-full" />
            <span className="text-gray-500">处理联系人信息</span>
          </div>
        </div>
      </div>
    </motion.div>
  );

  const renderCompleted = () => (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      className="text-center space-y-6 py-12"
    >
      <div className="flex justify-center">
        <CheckCircle className="w-16 h-16 text-green-500" />
      </div>

      <div>
        <h2 className="text-2xl font-bold text-white mb-3">数据导入完成</h2>
        <p className="text-gray-400">
          微信数据已成功导入到系统中
        </p>
      </div>

      {importResponse && (
        <div className="bg-gray-800/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700 max-w-md mx-auto">
          <div className="space-y-3 text-left">
            <div className="flex justify-between">
              <span className="text-gray-400">同步ID:</span>
              <span className="text-white">{importResponse.sync_id}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">状态:</span>
              <span className="text-green-400">{importResponse.status}</span>
            </div>
          </div>
        </div>
      )}

      {wechatStatus && (
        <div className="bg-gray-800/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700 max-w-md mx-auto">
          <h3 className="text-lg font-semibold text-white mb-4">数据统计</h3>
          <div className="space-y-3 text-left">
            <div className="flex justify-between">
              <span className="text-gray-400">消息数量:</span>
              <span className="text-white">{wechatStatus.messageCount.toLocaleString()} 条</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">联系人:</span>
              <span className="text-white">{wechatStatus.contactCount.toLocaleString()} 人</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">聊天会话:</span>
              <span className="text-white">{wechatStatus.chatCount.toLocaleString()} 个</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">最后同步:</span>
              <span className="text-white">
                {new Date(wechatStatus.lastSyncTime).toLocaleString('zh-CN')}
              </span>
            </div>
          </div>
        </div>
      )}

      <div className="flex justify-center gap-4">
        <button
          onClick={() => window.location.reload()}
          className="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
        >
          查看数据
        </button>
        <button
          onClick={onCancel}
          className="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
        >
          关闭
        </button>
      </div>
    </motion.div>
  );

  const renderError = () => (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      className="text-center space-y-6 py-12"
    >
      <div className="flex justify-center">
        <AlertCircle className="w-16 h-16 text-red-500" />
      </div>

      <div>
        <h2 className="text-2xl font-bold text-white mb-3">导入失败</h2>
        <p className="text-gray-400 mb-4">
          微信数据导入过程中遇到问题
        </p>
        {errorMessage && (
          <div className="bg-red-900/20 border border-red-700 rounded-lg p-4 text-red-400 text-sm max-w-md mx-auto">
            {errorMessage}
          </div>
        )}
      </div>

      <div className="flex justify-center gap-4">
        <button
          onClick={handleRetry}
          className="flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
        >
          <ArrowLeft className="w-4 h-4" />
          重新配置
        </button>
        <button
          onClick={onCancel}
          className="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
        >
          取消
        </button>
      </div>
    </motion.div>
  );

  const renderStepIndicator = () => (
    <div className="flex items-center justify-center mb-8">
      <div className="flex items-center gap-2">
        {steps.map((step, index) => {
          const isActive = steps.findIndex(s => s.id === currentStep) === index;
          const isCompleted = steps.findIndex(s => s.id === currentStep) > index;
          const isError = currentStep === 'error' && index === steps.length - 1;

          return (
            <React.Fragment key={step.id}>
              <div
                className={`
                  flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-all
                  ${isActive ? 'bg-blue-600 text-white' :
                    isCompleted ? 'bg-green-600/20 text-green-400' :
                    isError ? 'bg-red-600/20 text-red-400' :
                    'bg-gray-800 text-gray-400'}
                `}
              >
                {isCompleted ? (
                  <CheckCircle className="w-4 h-4" />
                ) : isError ? (
                  <AlertCircle className="w-4 h-4" />
                ) : (
                  <span className="w-4 h-4 flex items-center justify-center text-xs">
                    {index + 1}
                  </span>
                )}
                <span>{step.title}</span>
              </div>
              {index < steps.length - 1 && (
                <div className={`w-8 h-0.5 ${isCompleted ? 'bg-green-600' : 'bg-gray-700'}`} />
              )}
            </React.Fragment>
          );
        })}
      </div>
    </div>
  );

  return (
    <div className="min-h-screen bg-black text-white flex items-center justify-center p-6">
      <div className="w-full max-w-4xl">
        <div className="mb-8">
          {renderStepIndicator()}
        </div>

        <AnimatePresence mode="wait">
          {currentStep === 'source-select' && renderSourceSelect()}
          {currentStep === 'path-config' && renderPathConfig()}
          {currentStep === 'importing' && renderImporting()}
          {currentStep === 'completed' && renderCompleted()}
          {currentStep === 'error' && renderError()}
        </AnimatePresence>
      </div>
    </div>
  );
};

export default WeChatImportWizard;