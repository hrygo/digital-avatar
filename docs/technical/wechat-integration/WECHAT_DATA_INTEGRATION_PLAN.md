# 微信数据集成技术方案

<div align="center">

![Technical Design](https://img.shields.io/badge/Type-Technical%20Design-purple?style=flat-square)
![Version](https://img.shields.io/badge/Target-v0.4.0-blue?style=flat-square)
![Complexity](https://img.shields.io/badge/Complexity-High-orange?style=flat-square)

**真实微信数据读取、解析、分析与安全保护完整解决方案**

</div>

## 📋 方案概述

本文档详细描述了Digital Avatar v0.4.0版本中微信数据集成的完整技术方案，包括数据读取、解析、安全保护、AI分析和用户界面等各个层面的技术实现细节。

## 🎯 技术目标

### 核心能力
1. **真实数据读取**: 能够读取真实的微信聊天记录
2. **安全数据处理**: 确保用户隐私和数据安全
3. **智能分析增强**: 基于真实数据的AI分析能力
4. **用户友好体验**: 简单易用的数据导入和管理界面

### 技术约束
- **隐私保护**: 所有数据处理在本地完成，不上传云端
- **安全合规**: 符合相关法律法规要求
- **性能要求**: 支持大数据量处理，响应时间合理
- **兼容性**: 支持不同版本的微信数据格式

## 🏗️ 系统架构设计

### 整体架构图
```
┌─────────────────────────────────────────────────────────────────────┐
│                        Digital Avatar 系统架构                        │
├─────────────────────────────────────────────────────────────────────┤
│                           用户界面层 (Frontend)                        │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐     │
│  │ 数据导入向导 │ │ 仪表板界面  │ │ 分析结果展示 │ │ 隐私设置界面 │     │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘     │
├─────────────────────────────────────────────────────────────────────┤
│                           API服务层 (Backend API)                      │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐     │
│  │ 数据导入API  │ │ 数据查询API  │ │ 分析服务API  │ │ 系统管理API  │     │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘     │
├─────────────────────────────────────────────────────────────────────┤
│                           业务逻辑层 (Business Logic)                   │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐     │
│  │ 微信数据读取 │ │ 数据同步管理 │ │ 隐私脱敏处理 │ │ AI分析引擎   │     │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘     │
├─────────────────────────────────────────────────────────────────────┤
│                           数据存储层 (Data Storage)                     │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐     │
│  │ 临时缓存区   │ │ 加密数据库   │ │ 分析结果存储 │ │ 日志存储     │     │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘     │
├─────────────────────────────────────────────────────────────────────┤
│                           数据源层 (Data Sources)                      │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐     │
│  │ 微信数据库   │ │ 微信备份文件 │ │ 微信导出数据 │ │ 其他数据源   │     │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘     │
└─────────────────────────────────────────────────────────────────────┘
```

## 🔧 核心技术实现

### 1. 微信数据读取模块

#### 1.1 微信数据库解析器

**技术栈**: Go + SQLite3 + 加密解密

```go
// 微信数据库解析器结构
type WeChatDBParser struct {
    dbPath        string                 // 微信数据库路径
    db            *sql.DB               // 数据库连接
    encryptionKey []byte                // 解密密钥
    contacts      map[string]*Contact   // 联系人缓存
    messageCache  []*Message            // 消息缓存
}

// 微信消息结构
type Message struct {
    ID           int64                  `json:"id"`
    SvrID        int64                  `json:"svr_id"`
    CreateTime   int64                  `json:"create_time"`
    Talker       string                 `json:"talker"`
    Type         int                    `json:"type"`
    SubType      int                    `json:"sub_type"`
    IsSender     int                    `json:"is_sender"`
    Seq          int64                  `json:"seq"`
    Flag         int                    `json:"flag"`
    Status       int                    `json:"status"`
    Content      string                 `json:"content"`
    DisplayContent string               `json:"display_content"`
    LVBuffer     []byte                 `json:"lv_buffer"`
    TalkerID     int64                  `json:"talker_id"`
    BytesExtra   []byte                 `json:"bytes_extra"`
    CompressContent []byte              `json:"compress_content"`
}

// 数据库解析核心方法
func (p *WeChatDBParser) ParseDatabase() error {
    // 1. 打开微信数据库
    db, err := sql.Open("sqlite3", p.dbPath)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }
    defer db.Close()

    // 2. 读取联系人信息
    if err := p.parseContacts(db); err != nil {
        return fmt.Errorf("failed to parse contacts: %w", err)
    }

    // 3. 读取消息记录
    if err := p.parseMessages(db); err != nil {
        return fmt.Errorf("failed to parse messages: %w", err)
    }

    return nil
}

// 解析联系人信息
func (p *WeChatDBParser) parseContacts(db *sql.DB) error {
    query := `
        SELECT username, nickname, remark, type, headImgUrl
        FROM rcontact
        WHERE type & 1 = 1 AND username NOT LIKE 'gh_%'
    `

    rows, err := db.Query(query)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var contact Contact
        if err := rows.Scan(&contact.Username, &contact.Nickname,
                           &contact.Remark, &contact.Type, &contact.Avatar); err != nil {
            continue
        }
        p.contacts[contact.Username] = &contact
    }

    return nil
}

// 解析消息记录
func (p *WeChatDBParser) parseMessages(db *sql.DB) error {
    query := `
        SELECT msgId, msgSvrId, type, subType, isSender, createTime,
               talker, content, displayContent, status, lvbuffer
        FROM message
        ORDER BY createTime ASC
    `

    rows, err := db.Query(query)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var msg Message
        if err := rows.Scan(&msg.ID, &msg.SvrID, &msg.Type, &msg.SubType,
                           &msg.IsSender, &msg.CreateTime, &msg.Talker,
                           &msg.Content, &msg.DisplayContent, &msg.Status,
                           &msg.LVBuffer); err != nil {
            continue
        }

        // 解析消息内容
        if err := p.parseMessageContent(&msg); err == nil {
            p.messageCache = append(p.messageCache, &msg)
        }
    }

    return nil
}

// 解析消息内容（处理不同类型消息）
func (p *WeChatDBParser) parseMessageContent(msg *Message) error {
    switch msg.Type {
    case 1: // 文本消息
        // 文本消息直接使用content字段
        if msg.Content == "" {
            msg.Content = msg.DisplayContent
        }

    case 3: // 图片消息
        // 解析图片信息
        msg.DisplayContent = "[图片]"

    case 34: // 语音消息
        msg.DisplayContent = "[语音]"

    case 43: // 视频消息
        msg.DisplayContent = "[视频]"

    case 49: // 文件消息/链接消息
        // 解析文件或链接信息
        if strings.Contains(msg.Content, "msg") {
            msg.DisplayContent = "[文件]"
        } else {
            msg.DisplayContent = "[链接]"
        }

    default:
        msg.DisplayContent = "[未知消息类型]"
    }

    return nil
}
```

#### 1.2 微信备份文件解析器

```go
// 微信备份文件解析器
type WeChatBackupParser struct {
    backupPath   string
    tempDir      string
    dbParser     *WeChatDBParser
}

// 解析微信备份文件
func (p *WeChatBackupParser) ParseBackup() error {
    // 1. 解压备份文件
    if err := p.extractBackup(); err != nil {
        return fmt.Errorf("failed to extract backup: %w", err)
    }

    // 2. 查找微信数据库文件
    dbPath, err := p.findWeChatDB()
    if err != nil {
        return fmt.Errorf("failed to find WeChat DB: %w", err)
    }

    // 3. 使用数据库解析器解析
    p.dbParser = &WeChatDBParser{dbPath: dbPath}
    return p.dbParser.ParseDatabase()
}

// 解压备份文件
func (p *WeChatBackupParser) extractBackup() error {
    // 创建临时目录
    tempDir, err := ioutil.TempDir("", "wechat_backup_")
    if err != nil {
        return err
    }
    p.tempDir = tempDir

    // 根据备份文件类型选择解压方式
    if strings.HasSuffix(p.backupPath, ".tar") {
        return p.extractTar(p.backupPath, tempDir)
    } else if strings.HasSuffix(p.backupPath, ".zip") {
        return p.extractZip(p.backupPath, tempDir)
    }

    return fmt.Errorf("unsupported backup format")
}
```

### 2. 数据同步机制

#### 2.1 增量同步服务

```go
// 数据同步服务
type DataSyncService struct {
    parser       *WeChatDBParser
    lastSyncTime int64
    syncStatus   SyncStatus
    mu           sync.RWMutex
}

// 同步状态
type SyncStatus struct {
    Status       string    `json:"status"`        // idle, syncing, completed, error
    Progress     float64   `json:"progress"`      // 0.0 - 1.0
    LastSyncTime int64     `json:"last_sync_time"`
    TotalCount   int       `json:"total_count"`
    ProcessedCount int     `json:"processed_count"`
    ErrorMsg     string    `json:"error_msg,omitempty"`
}

// 执行增量同步
func (s *DataSyncService) IncrementalSync() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    s.syncStatus.Status = "syncing"
    s.syncStatus.Progress = 0.0

    // 获取最后同步时间后的新消息
    newMessages, err := s.getMessagesSince(s.lastSyncTime)
    if err != nil {
        s.syncStatus.Status = "error"
        s.syncStatus.ErrorMsg = err.Error()
        return err
    }

    // 处理新消息
    total := len(newMessages)
    for i, msg := range newMessages {
        if err := s.processMessage(msg); err != nil {
            log.Printf("Failed to process message %d: %v", msg.ID, err)
            continue
        }

        s.syncStatus.ProcessedCount = i + 1
        s.syncStatus.Progress = float64(i+1) / float64(total)
    }

    // 更新同步时间
    s.lastSyncTime = time.Now().Unix()
    s.syncStatus.Status = "completed"
    s.syncStatus.LastSyncTime = s.lastSyncTime

    return nil
}

// 获取指定时间后的消息
func (s *DataSyncService) getMessagesSince(timestamp int64) ([]*Message, error) {
    query := `
        SELECT msgId, msgSvrId, type, subType, isSender, createTime,
               talker, content, displayContent, status
        FROM message
        WHERE createTime > ?
        ORDER BY createTime ASC
    `

    rows, err := s.parser.db.Query(query, timestamp)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []*Message
    for rows.Next() {
        var msg Message
        if err := rows.Scan(&msg.ID, &msg.SvrID, &msg.Type, &msg.SubType,
                           &msg.IsSender, &msg.CreateTime, &msg.Talker,
                           &msg.Content, &msg.DisplayContent, &msg.Status); err != nil {
            continue
        }

        s.parser.parseMessageContent(&msg)
        messages = append(messages, &msg)
    }

    return messages, nil
}
```

#### 2.2 前端同步服务

```typescript
// 前端数据同步服务
export class WeChatSyncService {
    private apiClient: ApiClient;
    private syncStatus$: BehaviorSubject<SyncStatus>;

    constructor() {
        this.apiClient = new ApiClient();
        this.syncStatus$ = new BehaviorSubject<SyncStatus>({
            status: 'idle',
            progress: 0,
            lastSyncTime: 0,
            totalCount: 0,
            processedCount: 0
        });
    }

    // 开始全量同步
    async startFullSync(): Promise<void> {
        try {
            await this.apiClient.post('/api/v1/wechat/sync/full');
            await this.monitorSyncProgress();
        } catch (error) {
            throw new Error(`Full sync failed: ${error.message}`);
        }
    }

    // 开始增量同步
    async startIncrementalSync(): Promise<void> {
        try {
            await this.apiClient.post('/api/v1/wechat/sync/incremental');
            await this.monitorSyncProgress();
        } catch (error) {
            throw new Error(`Incremental sync failed: ${error.message}`);
        }
    }

    // 监控同步进度
    private async monitorSyncProgress(): Promise<void> {
        const checkProgress = async () => {
            const status = await this.apiClient.get<SyncStatus>('/api/v1/wechat/sync/status');
            this.syncStatus$.next(status);

            if (status.status === 'syncing') {
                setTimeout(checkProgress, 1000); // 每秒检查一次
            }
        };

        await checkProgress();
    }

    // 获取同步状态
    getSyncStatus(): Observable<SyncStatus> {
        return this.syncStatus$.asObservable();
    }
}
```

### 3. 隐私保护与数据脱敏

#### 3.1 增强型PII检测器

```go
// 增强型PII检测器
type EnhancedPIIDetector struct {
    patterns map[string]*regexp.Regexp
    names    map[string]bool // 常见中文姓名
}

// 初始化PII检测器
func NewEnhancedPIIDetector() *EnhancedPIIDetector {
    detector := &EnhancedPIIDetector{
        patterns: make(map[string]*regexp.Regexp),
        names:    make(map[string]bool),
    }

    // 初始化检测模式
    detector.initPatterns()
    detector.initChineseNames()

    return detector
}

// 初始化检测模式
func (d *EnhancedPIIDetector) initPatterns() {
    // 手机号模式
    d.patterns["phone"] = regexp.MustCompile(`1[3-9]\d{9}`)

    // 邮箱模式
    d.patterns["email"] = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

    // 身份证模式
    d.patterns["id_card"] = regexp.MustCompile(`\d{17}[\dXx]`)

    // 银行卡模式
    d.patterns["bank_card"] = regexp.MustCompile(`\d{16,19}`)

    // 地址模式（简化版）
    d.patterns["address"] = regexp.MustCompile(`[省市区县].*[街道路].*[号楼]?`)

    // 微信号模式
    d.patterns["wechat_id"] = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9_-]{5,19}`)

    // QQ号模式
    d.patterns["qq"] = regexp.MustCompile(`[1-9]\d{4,10}`)
}

// 检测消息中的PII信息
func (d *EnhancedPIIDetector) DetectPII(content string) []PIIInfo {
    var piiInfos []PIIInfo

    // 检测各种PII类型
    for piiType, pattern := range d.patterns {
        matches := pattern.FindAllStringSubmatch(content, -1)
        for _, match := range matches {
            piiInfos = append(piiInfos, PIIInfo{
                Type:   piiType,
                Value:  match[0],
                Start:  strings.Index(content, match[0]),
                End:    strings.Index(content, match[0]) + len(match[0]),
                Confidence: d.calculateConfidence(piiType, match[0]),
            })
        }
    }

    // 检测中文姓名
    d.detectChineseNames(content, &piiInfos)

    return piiInfos
}

// 检测中文姓名
func (d *EnhancedPIIDetector) detectChineseNames(content string, piiInfos *[]PIIInfo) {
    // 使用简单的2-3字中文名检测
    namePattern := regexp.MustCompile(`[\p{Han}]{2,3}`)
    matches := namePattern.FindAllString(content, -1)

    for _, name := range matches {
        if d.isLikelyChineseName(name) {
            *piiInfos = append(*piiInfos, PIIInfo{
                Type:   "chinese_name",
                Value:  name,
                Start:  strings.Index(content, name),
                End:    strings.Index(content, name) + len(name),
                Confidence: 0.7,
            })
        }
    }
}

// 判断是否可能是中文姓名
func (d *EnhancedPIIDetector) isLikelyChineseName(name string) bool {
    // 常见姓氏判断
    commonSurnames := []string{"王", "李", "张", "刘", "陈", "杨", "赵", "黄", "周", "吴"}

    if len(name) < 2 || len(name) > 3 {
        return false
    }

    firstChar := string(name[0])
    for _, surname := range commonSurnames {
        if firstChar == surname {
            return true
        }
    }

    return false
}

// 计算置信度
func (d *EnhancedPIIDetector) calculateConfidence(piiType, value string) float64 {
    switch piiType {
    case "phone":
        return 0.95
    case "email":
        return 0.90
    case "id_card":
        return 0.98
    case "bank_card":
        return 0.85
    case "wechat_id":
        return 0.80
    default:
        return 0.70
    }
}
```

#### 3.2 脱敏处理器

```go
// 脱敏处理器
type Anonymizer struct {
    detector *EnhancedPIIDetector
    salt     []byte
}

// 初始化脱敏器
func NewAnonymizer() *Anonymizer {
    salt := make([]byte, 32)
    rand.Read(salt)

    return &Anonymizer{
        detector: NewEnhancedPIIDetector(),
        salt:     salt,
    }
}

// 脱敏处理
func (a *Anonymizer) AnonymizeText(content string) string {
    piiInfos := a.detector.DetectPII(content)

    // 按位置倒序排列，避免替换影响位置
    sort.Slice(piiInfos, func(i, j int) bool {
        return piiInfos[i].Start > piiInfos[j].Start
    })

    result := content
    for _, pii := range piiInfos {
        replacement := a.generateReplacement(pii)
        result = result[:pii.Start] + replacement + result[pii.End:]
    }

    return result
}

// 生成替换内容
func (a *Anonymizer) generateReplacement(pii PIIInfo) string {
    hash := sha256.Sum256(append([]byte(pii.Value), a.salt...))
    hashStr := hex.EncodeToString(hash[:])[:8]

    switch pii.Type {
    case "phone":
        return fmt.Sprintf("手机号[%s]", hashStr)
    case "email":
        return fmt.Sprintf("邮箱[%s]", hashStr)
    case "id_card":
        return fmt.Sprintf("身份证[%s]", hashStr)
    case "chinese_name":
        return fmt.Sprintf("姓名[%s]", hashStr)
    case "wechat_id":
        return fmt.Sprintf("微信号[%s]", hashStr)
    case "address":
        return fmt.Sprintf("地址[%s]", hashStr)
    default:
        return fmt.Sprintf("敏感信息[%s]", hashStr)
    }
}
```

### 4. 前端界面组件

#### 4.1 微信数据导入组件

```typescript
// 微信数据导入组件
import React, { useState, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

interface ImportWizardProps {
  onImportComplete: (result: ImportResult) => void;
}

const WeChatImportWizard: React.FC<ImportWizardProps> = ({ onImportComplete }) => {
  const [currentStep, setCurrentStep] = useState(0);
  const [importType, setImportType] = useState<'database' | 'backup' | 'export'>('database');
  const [importPath, setImportPath] = useState('');
  const [isImporting, setIsImporting] = useState(false);
  const [importProgress, setImportProgress] = useState(0);
  const [importResult, setImportResult] = useState<ImportResult | null>(null);

  const steps = [
    { title: '选择导入方式', component: Step1_SelectType },
    { title: '选择数据源', component: Step2_SelectSource },
    { title: '数据预览', component: Step3_PreviewData },
    { title: '导入执行', component: Step4_Importing },
    { title: '完成', component: Step5_Completed },
  ];

  const CurrentStepComponent = steps[currentStep].component;

  const handleNext = useCallback(() => {
    if (currentStep < steps.length - 1) {
      setCurrentStep(currentStep + 1);
    }
  }, [currentStep, steps.length]);

  const handlePrevious = useCallback(() => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  }, [currentStep]);

  const handleImport = useCallback(async () => {
    setIsImporting(true);
    setImportProgress(0);

    try {
      const result = await WeChatSyncService.importData({
        type: importType,
        path: importPath,
        onProgress: (progress) => setImportProgress(progress),
      });

      setImportResult(result);
      onImportComplete(result);
      setCurrentStep(steps.length - 1);
    } catch (error) {
      console.error('Import failed:', error);
    } finally {
      setIsImporting(false);
    }
  }, [importType, importPath, steps.length, onImportComplete]);

  return (
    <motion.div
      className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
    >
      <motion.div
        className="bg-slate-800/90 backdrop-blur-xl rounded-2xl p-8 max-w-2xl w-full max-h-[80vh] overflow-y-auto border border-slate-700/50"
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.9, opacity: 0 }}
      >
        {/* 步骤指示器 */}
        <div className="mb-8">
          <div className="flex items-center justify-between mb-4">
            {steps.map((step, index) => (
              <div key={index} className="flex items-center">
                <div
                  className={`w-10 h-10 rounded-full flex items-center justify-center text-sm font-medium transition-colors ${
                    index <= currentStep
                      ? 'bg-blue-500 text-white'
                      : 'bg-slate-700 text-slate-400'
                  }`}
                >
                  {index + 1}
                </div>
                {index < steps.length - 1 && (
                  <div
                    className={`h-1 w-20 mx-2 transition-colors ${
                      index < currentStep ? 'bg-blue-500' : 'bg-slate-700'
                    }`}
                  />
                )}
              </div>
            ))}
          </div>
          <h3 className="text-xl font-semibold text-white text-center">
            {steps[currentStep].title}
          </h3>
        </div>

        {/* 步骤内容 */}
        <AnimatePresence mode="wait">
          <motion.div
            key={currentStep}
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -20 }}
            transition={{ duration: 0.2 }}
          >
            <CurrentStepComponent
              importType={importType}
              onImportTypeChange={setImportType}
              importPath={importPath}
              onImportPathChange={setImportPath}
              onImport={handleImport}
              isImporting={isImporting}
              importProgress={importProgress}
              importResult={importResult}
            />
          </motion.div>
        </AnimatePresence>

        {/* 导航按钮 */}
        <div className="flex justify-between mt-8">
          <button
            onClick={handlePrevious}
            disabled={currentStep === 0}
            className="px-6 py-2 bg-slate-700 text-white rounded-lg hover:bg-slate-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            上一步
          </button>
          <button
            onClick={handleNext}
            disabled={currentStep === steps.length - 1}
            className="px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            下一步
          </button>
        </div>
      </motion.div>
    </motion.div>
  );
};

// 步骤1：选择导入类型
const Step1_SelectType: React.FC<StepProps> = ({ importType, onImportTypeChange, onNext }) => (
  <div className="space-y-6">
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      {[
        { type: 'database', title: '微信数据库', desc: '直接读取微信客户端数据库' },
        { type: 'backup', title: '微信备份', desc: '从微信备份文件导入' },
        { type: 'export', title: '导出文件', desc: '从微信导出的文件导入' },
      ].map((option) => (
        <motion.button
          key={option.type}
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          onClick={() => onImportTypeChange(option.type as any)}
          className={`p-6 rounded-xl border-2 transition-all ${
            importType === option.type
              ? 'border-blue-500 bg-blue-500/10'
              : 'border-slate-600 bg-slate-800/50 hover:border-slate-500'
          }`}
        >
          <h4 className="text-lg font-semibold text-white mb-2">{option.title}</h4>
          <p className="text-slate-400 text-sm">{option.desc}</p>
        </motion.button>
      ))}
    </div>
  </div>
);

// 其他步骤组件...
const Step2_SelectSource: React.FC<StepProps> = ({ importPath, onImportPathChange, onNext }) => (
  <div className="space-y-6">
    <div>
      <label className="block text-sm font-medium text-slate-300 mb-2">
        选择微信数据路径
      </label>
      <div className="flex gap-4">
        <input
          type="text"
          value={importPath}
          onChange={(e) => onImportPathChange(e.target.value)}
          placeholder="请选择微信数据文件或文件夹路径"
          className="flex-1 px-4 py-2 bg-slate-700/50 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <button className="px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors">
          浏览
        </button>
      </div>
    </div>
  </div>
);

// 继续定义其他步骤组件...

export default WeChatImportWizard;
```

#### 4.2 数据同步状态组件

```typescript
// 数据同步状态组件
import React, { useEffect, useState } from 'react';
import { motion } from 'framer-motion';

const SyncStatusIndicator: React.FC = () => {
  const [syncStatus, setSyncStatus] = useState<SyncStatus>({
    status: 'idle',
    progress: 0,
    lastSyncTime: 0,
    totalCount: 0,
    processedCount: 0,
  });

  const [wechatStatus, setWeChatStatus] = useState<WeChatStatus>({
    connected: false,
    dbPath: '',
    lastScanTime: 0,
    messageCount: 0,
    contactCount: 0,
  });

  useEffect(() => {
    const syncService = new WeChatSyncService();
    const subscription = syncService.getSyncStatus().subscribe(setSyncStatus);

    return () => subscription.unsubscribe();
  }, []);

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed': return 'text-green-400';
      case 'syncing': return 'text-blue-400';
      case 'error': return 'text-red-400';
      default: return 'text-slate-400';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'idle': return '未同步';
      case 'syncing': return '同步中';
      case 'completed': return '已同步';
      case 'error': return '同步失败';
      default: return '未知状态';
    }
  };

  return (
    <motion.div
      className="bg-slate-800/50 backdrop-blur-sm rounded-xl p-6 border border-slate-700/50"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
    >
      <h3 className="text-lg font-semibold text-white mb-4">系统状态</h3>

      {/* 微信数据库状态 */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-3">
          <span className="text-slate-300">微信数据库</span>
          <span className={`text-sm font-medium ${getStatusColor(wechatStatus.connected ? 'completed' : 'error')}`}>
            {wechatStatus.connected ? '已连接' : '未连接'}
          </span>
        </div>

        {wechatStatus.connected && (
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div className="bg-slate-700/30 rounded-lg p-3">
              <div className="text-slate-400">消息数量</div>
              <div className="text-xl font-semibold text-white">
                {wechatStatus.messageCount.toLocaleString()}
              </div>
            </div>
            <div className="bg-slate-700/30 rounded-lg p-3">
              <div className="text-slate-400">联系人数量</div>
              <div className="text-xl font-semibold text-white">
                {wechatStatus.contactCount.toLocaleString()}
              </div>
            </div>
          </div>
        )}
      </div>

      {/* 同步状态 */}
      <div>
        <div className="flex items-center justify-between mb-3">
          <span className="text-slate-300">数据同步</span>
          <span className={`text-sm font-medium ${getStatusColor(syncStatus.status)}`}>
            {getStatusText(syncStatus.status)}
          </span>
        </div>

        {syncStatus.status === 'syncing' && (
          <div className="mb-3">
            <div className="flex justify-between text-sm text-slate-400 mb-1">
              <span>进度</span>
              <span>{Math.round(syncStatus.progress * 100)}%</span>
            </div>
            <div className="w-full bg-slate-700 rounded-full h-2">
              <motion.div
                className="bg-blue-500 h-2 rounded-full"
                style={{ width: `${syncStatus.progress * 100}%` }}
                initial={{ width: 0 }}
                animate={{ width: `${syncStatus.progress * 100}%` }}
                transition={{ duration: 0.3 }}
              />
            </div>
            <div className="text-sm text-slate-400 mt-1">
              {syncStatus.processedCount} / {syncStatus.totalCount}
            </div>
          </div>
        )}

        {syncStatus.lastSyncTime > 0 && (
          <div className="text-sm text-slate-400">
            最后同步: {new Date(syncStatus.lastSyncTime * 1000).toLocaleString()}
          </div>
        )}
      </div>

      {/* 操作按钮 */}
      <div className="flex gap-3 mt-6">
        <button
          onClick={() => {/* 开始全量同步 */}}
          disabled={syncStatus.status === 'syncing'}
          className="flex-1 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          全量同步
        </button>
        <button
          onClick={() => {/* 开始增量同步 */}}
          disabled={syncStatus.status === 'syncing' || syncStatus.lastSyncTime === 0}
          className="flex-1 px-4 py-2 bg-slate-700 text-white rounded-lg hover:bg-slate-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          增量同步
        </button>
      </div>
    </motion.div>
  );
};

export default SyncStatusIndicator;
```

## 🛡️ 安全架构设计

### 1. 数据加密方案

```go
// 加密管理器
type CryptoManager struct {
    masterKey []byte
}

// 生成主密钥
func (c *CryptoManager) generateMasterKey() error {
    c.masterKey = make([]byte, 32)
    _, err := rand.Read(c.masterKey)
    return err
}

// 加密敏感数据
func (c *CryptoManager) Encrypt(plaintext string) (string, error) {
    block, err := aes.NewCipher(c.masterKey)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// 解密敏感数据
func (c *CryptoManager) Decrypt(ciphertext string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(c.masterKey)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return "", errors.New("ciphertext too short")
    }

    nonce, ciphertext_bytes := data[:nonceSize], data[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext_bytes, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}
```

### 2. 访问控制

```go
// 访问控制管理器
type AccessControlManager struct {
    permissions map[string][]string
    sessions    map[string]*Session
}

// 会话结构
type Session struct {
    ID        string
    UserID    string
    CreatedAt time.Time
    ExpiresAt time.Time
    IsActive  bool
}

// 检查权限
func (a *AccessControlManager) CheckPermission(sessionID, resource, action string) bool {
    session, exists := a.sessions[sessionID]
    if !exists || !session.IsActive || session.ExpiresAt.Before(time.Now()) {
        return false
    }

    userPermissions, exists := a.permissions[session.UserID]
    if !exists {
        return false
    }

    requiredPermission := fmt.Sprintf("%s:%s", resource, action)
    for _, permission := range userPermissions {
        if permission == requiredPermission || permission == "*" {
            return true
        }
    }

    return false
}
```

## 📊 性能优化策略

### 1. 数据库优化

```sql
-- 创建优化的索引
CREATE INDEX idx_message_create_time ON message(create_time);
CREATE INDEX idx_message_talker ON message(talker);
CREATE INDEX idx_message_type ON message(type);
CREATE INDEX idx_contact_username ON rcontact(username);

-- 创建分区表（如果数据量很大）
CREATE TABLE message_partitioned (
    LIKE message INCLUDING ALL
) PARTITION BY RANGE (create_time);

-- 按月分区
CREATE TABLE message_2025_11 PARTITION OF message_partitioned
    FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');
```

### 2. 缓存策略

```go
// 缓存管理器
type CacheManager struct {
    cache      map[string]interface{}
    expiration map[string]time.Time
    mu         sync.RWMutex
}

// 设置缓存
func (c *CacheManager) Set(key string, value interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.cache[key] = value
    c.expiration[key] = time.Now().Add(ttl)
}

// 获取缓存
func (c *CacheManager) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    if expiration, exists := c.expiration[key]; exists {
        if time.Now().Before(expiration) {
            return c.cache[key], true
        }
        // 过期，清理
        delete(c.cache, key)
        delete(c.expiration, key)
    }

    return nil, false
}
```

## 🧪 测试方案

### 1. 单元测试

```go
// PII检测器测试
func TestPIIDetector(t *testing.T) {
    detector := NewEnhancedPIIDetector()

    testCases := []struct {
        content string
        expected int
    }{
        {"我的手机号是13812345678", 1},
        {"邮箱地址是test@example.com", 1},
        {"身份证号是110101199001011234", 1},
        {"联系王先生和李女士", 2},
    }

    for _, tc := range testCases {
        piiInfos := detector.DetectPII(tc.content)
        if len(piiInfos) != tc.expected {
            t.Errorf("Expected %d PII infos, got %d", tc.expected, len(piiInfos))
        }
    }
}
```

### 2. 集成测试

```go
// 数据导入集成测试
func TestWeChatDataImport(t *testing.T) {
    // 设置测试环境
    testDB := setupTestDatabase(t)
    defer testDB.Close()

    // 创建测试数据
    testMessages := createTestMessages(t, testDB)

    // 执行导入
    parser := &WeChatDBParser{dbPath: testDB.Path()}
    err := parser.ParseDatabase()
    if err != nil {
        t.Fatalf("Failed to parse database: %v", err)
    }

    // 验证结果
    if len(parser.messageCache) != len(testMessages) {
        t.Errorf("Expected %d messages, got %d", len(testMessages), len(parser.messageCache))
    }
}
```

## 📅 实施计划

### Phase 1: 基础架构 (第1-2天)
- [ ] 搭建微信数据解析框架
- [ ] 实现基础数据库连接和读取
- [ ] 创建数据模型定义
- [ ] 设置开发测试环境

### Phase 2: 数据读取 (第3-5天)
- [ ] 实现微信数据库解析器
- [ ] 开发备份文件解析功能
- [ ] 创建数据验证机制
- [ ] 实现基础错误处理

### Phase 3: 安全保护 (第6-7天)
- [ ] 集成增强型PII检测器
- [ ] 实现数据脱敏功能
- [ ] 添加加密存储机制
- [ ] 完善访问控制

### Phase 4: 前端界面 (第8-10天)
- [ ] 开发数据导入向导
- [ ] 实现同步状态监控
- [ ] 创建数据预览功能
- [ ] 优化用户体验

### Phase 5: 测试与优化 (第11-14天)
- [ ] 编写完整的测试用例
- [ ] 性能优化和压力测试
- [ ] 安全测试和漏洞修复
- [ ] 文档编写和代码审查

## 📚 相关文档

- [微信数据结构分析](../analysis/wechat-database-structure.md)
- [隐私保护技术方案](../privacy/pii-protection.md)
- [API接口文档](../../api/wechat-api.md)
- [前端组件库文档](../../frontend/components.md)

---

**文档版本**: v1.0
**创建时间**: 2025年11月19日
**最后更新**: 2025年11月19日
**文档状态**: 设计完成，待实施

---

*本文档详细描述了微信数据集成的完整技术方案，为开发团队提供了明确的实施指导。*